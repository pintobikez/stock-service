package middleware

import (
	"crypto/tls"
	"github.com/labstack/echo"
	strut "github.com/pintobikez/stock-service/api/structures"
	cnfs "github.com/pintobikez/stock-service/config/structures"
	"net/http"
	"regexp"
)

// Authorization Middleware
func Authorization(authConfig *cnfs.AuthConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//connect to an authentication service and do authentication
			req, err := http.NewRequest("POST", authConfig.Url, nil)
			if err != nil {
				return c.JSON(http.StatusServiceUnavailable, &strut.ErrResponse{strut.ErrContent{http.StatusServiceUnavailable, "Authorization service is down"}})
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
				return c.JSON(http.StatusServiceUnavailable, &strut.ErrResponse{strut.ErrContent{http.StatusServiceUnavailable, "Authorization service is down"}})
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				return c.JSON(http.StatusUnauthorized, &strut.ErrResponse{strut.ErrContent{res.StatusCode, "Authorization error"}})
			}

			return next(c)
		}
	}
}
