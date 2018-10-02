package main

import (
	"context"
	"crypto/tls"
	middleware "github.com/dafiti/echo-middleware"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/color"
	"github.com/labstack/gommon/log"
	api "github.com/pintobikez/stock-service/api"
	uti "github.com/pintobikez/stock-service/config"
	cnfs "github.com/pintobikez/stock-service/config/structures"
	lg "github.com/pintobikez/stock-service/log"
	pub "github.com/pintobikez/stock-service/publisher"
	pb "github.com/pintobikez/stock-service/publisher/rabbitmq"
	rep "github.com/pintobikez/stock-service/repository"
	mysql "github.com/pintobikez/stock-service/repository/mysql"
	srv "github.com/pintobikez/stock-service/server"
	"gopkg.in/urfave/cli.v1"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"time"
)

var (
	repo      rep.Repository
	pubsub    pub.PubSub
	apiStruct *api.API
)

// Start Http Server
func Handler(c *cli.Context) error {

	// Echo instance
	e := &srv.Server{echo.New()}
	e.HTTPErrorHandler = api.ServerErrorHandler
	e.Logger.SetLevel(log.INFO)
	e.Logger.SetOutput(lg.File(c.String("log-folder") + "/app.log"))

	// Middlewares
	e.Use(middleware.LoggerWithOutput(lg.File(c.String("log-folder") + "/access.log")))

	if c.String("newrelic-appname") != "" && c.String("newrelic-license-key") != "" {
		e.Use(middleware.NewRelic(
			c.String("newrelic-appname"),
			c.String("newrelic-license-key"),
		))
	}

	e.Use(mw.Recover())
	e.Use(mw.Secure())
	e.Use(mw.RequestID())
	e.Pre(mw.RemoveTrailingSlash())

	//loads Authorization service
	authConfig := new(cnfs.AuthConfig)
	err := uti.LoadConfigFile(c.String("auth-file"), authConfig)
	if err != nil {
		e.Logger.Infof("Error in Authorization service %s", err.Error())
	} else {
		e.Use(Authorization(authConfig))
	}

	//loads db connection
	dbConfig := new(cnfs.DatabaseConfig)
	err = uti.LoadConfigFile(c.String("database-file"), dbConfig)
	if err != nil {
		e.Logger.Fatal(err)
	}

	repo, err = mysql.New(dbConfig)
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Database connect
	err = repo.Connect()
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer repo.Disconnect()

	//loads rabbitmq config
	rbcnfg := new(cnfs.PublisherConfig)
	err = uti.LoadConfigFile(c.String("publisher-file"), rbcnfg)
	if err != nil {
		e.Logger.Fatal(err)
	}

	// RabbitMQ connect
	pubsub, err = pb.New(rbcnfg)
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer pubsub.Close()

	apiStruct = api.New(repo, pubsub)

	// Routes => api
	e.GET("/health", apiStruct.HealthStatus(), mw.CORSWithConfig(
		mw.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{echo.GET, echo.OPTIONS, echo.HEAD},
		},
	))

	e.PUT("/stock/:sku", apiStruct.PutStock(), mw.CORSWithConfig(
		mw.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{echo.PUT, echo.OPTIONS, echo.HEAD},
		},
	))
	e.PUT("/reservation/:sku", apiStruct.PutReservation(), mw.CORSWithConfig(
		mw.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{echo.PUT, echo.OPTIONS, echo.HEAD},
		},
	))
	e.DELETE("/reservation/:sku", apiStruct.RemoveReservation(), mw.CORSWithConfig(
		mw.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{echo.DELETE, echo.OPTIONS, echo.HEAD},
		},
	))
	e.GET("/stock/:sku", apiStruct.GetStock(), mw.CORSWithConfig(
		mw.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{echo.GET, echo.OPTIONS, echo.HEAD},
		},
	))

	if c.String("revision-file") != "" {
		e.File("/rev.txt", c.String("revision-file"))
	}

	if swagger := c.String("swagger-file"); swagger != "" {
		g := e.Group("/docs")
		g.Use(mw.CORSWithConfig(
			mw.CORSConfig{
				AllowOrigins: []string{"http://petstore.swagger.io"},
				AllowMethods: []string{echo.GET, echo.HEAD},
			},
		))

		g.GET("", func(c echo.Context) error {
			return c.File(swagger)
		})
	}

	// Start server
	colorer := color.New()
	colorer.Printf("⇛ %s service - %s\n", appName, color.Green(version))

	go func() {
		if err := start(e, c); err != nil {
			colorer.Printf(color.Red("⇛ shutting down the server\n"))
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	return nil
}

// Start http server
func start(e *srv.Server, c *cli.Context) error {

	if c.String("ssl-cert") != "" && c.String("ssl-key") != "" {
		return e.StartTLS(
			c.String("listen"),
			c.String("ssl-cert"),
			c.String("ssl-key"),
		)
	}

	return e.Start(c.String("listen"))
}

// Authorization Middleware
func Authorization(authConfig *cnfs.AuthConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//connect to an authentication service and do authentication
			req, err := http.NewRequest("POST", authConfig.Url, nil)
			if err != nil {
				return c.JSON(http.StatusServiceUnavailable, &api.ErrResponse{api.ErrContent{http.StatusServiceUnavailable, "Authorization service is down"}})
			}
			if len(authConfig.Headers) > 0 {
				for k, v := range authConfig.Headers {
					req.Header.Set(k, v)
				}
			}
			req.Header.Set("Authorization", c.Request().Header.Get(echo.HeaderAuthorization))
			req.Close = true

			// check if it is an https request
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: regexp.MustCompile("^https://").MatchString(authConfig.Url)},
			}

			client := &http.Client{Transport: tr}
			res, err := client.Do(req)
			if err != nil {
				return c.JSON(http.StatusServiceUnavailable, &api.ErrResponse{api.ErrContent{http.StatusServiceUnavailable, "Authorization service is down"}})
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				return c.JSON(http.StatusUnauthorized, &api.ErrResponse{api.ErrContent{res.StatusCode, "Authorization error"}})
			}

			return next(c)
		}
	}
}
