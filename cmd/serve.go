package main

import (
	api "bitbucket.org/ricardomvpinto/stock-service/api"
	cnf "bitbucket.org/ricardomvpinto/stock-service/config"
	lg "bitbucket.org/ricardomvpinto/stock-service/log"
	pub "bitbucket.org/ricardomvpinto/stock-service/publisher"
	rep "bitbucket.org/ricardomvpinto/stock-service/repository"
	srv "bitbucket.org/ricardomvpinto/stock-service/server"
	"context"
	middleware "github.com/dafiti/echo-middleware"
	inst "github.com/dafiti/go-instrument"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/color"
	"github.com/labstack/gommon/log"
	"gopkg.in/urfave/cli.v1"
	"os"
	"os/signal"
	"strconv"
	"time"
)

var (
	instrument inst.Instrument    = new(inst.Dummy)
	repo       *rep.Repository    = new(rep.Repository)
	pubsub     *pub.FilePublisher = new(pub.FilePublisher)
)

// Start Http Server
func Serve(c *cli.Context) error {

	// Echo instance
	e := &srv.Server{echo.New()}
	e.HTTPErrorHandler = api.Error
	e.Logger.SetLevel(log.INFO)
	e.Logger.SetOutput(lg.File(c.String("log-folder") + "/app.log"))

	// Middlewares
	e.Use(middleware.LoggerWithOutput(lg.File(c.String("log-folder") + "/access.log")))

	if c.String("newrelic-appname") != "" && c.String("newrelic-license-key") != "" {
		e.Use(middleware.NewRelic(
			c.String("newrelic-appname"),
			c.String("newrelic-license-key"),
		))

		instrument = new(inst.NewRelic)
	}

	e.Use(mw.Recover())
	e.Use(mw.Secure())
	e.Use(mw.RequestID())
	e.Pre(mw.RemoveTrailingSlash())

	//loads db connection
	err, stringConn := buildStringConnection(c.String("database-file"))
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Database connect
	err = repo.ConnectDB(stringConn)
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Routes => api
	e.PUT("/stock/:sku", api.PutStock(repo, pubsub))
	e.PUT("/reservation/:sku", api.PutReservation(repo, pubsub))
	e.DELETE("/reservation/:sku", api.RemoveReservation(repo, pubsub))
	e.GET("/stock/:sku", api.GetStock(repo))

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
	defer repo.DisconnectDB()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	return nil
}

// Start http server
func start(e *srv.Server, c *cli.Context) error {
	return e.Start(c.String("listen"))
}

func buildStringConnection(filename string) (error, string) {
	t, err := cnf.LoadConfigFile(filename)
	if err != nil {
		return err, ""
	}
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	stringConn := t.Driver.User + ":" + t.Driver.Pw
	stringConn += "@tcp(" + t.Driver.Host + ":" + strconv.Itoa(t.Driver.Port) + ")"
	stringConn += "/" + t.Driver.Schema + "?charset=utf8"

	return nil, stringConn
}
