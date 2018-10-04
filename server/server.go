package server

import (
	"github.com/coreos/go-systemd/activation"
	"github.com/labstack/echo"
	strut "github.com/pintobikez/stock-service/api/structures"
	"net/http"
)

type Server struct {
	*echo.Echo
}

// Start starts an HTTP server.
func (srv *Server) Start(address string) error {

	listeners, err := activation.Listeners(true)

	if err != nil {
		return err
	}

	if len(listeners) > 0 {
		srv.Echo.Listener = listeners[0]
	}

	return srv.Echo.Start(address)
}

// ServerErrorHandler sets the format of the error to be return by the server
func ServerErrorHandler(err error, c echo.Context) {
	code := http.StatusServiceUnavailable
	msg := http.StatusText(code)

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message.(string)
	}

	if c.Echo().Debug {
		msg = err.Error()
	}

	content := map[string]interface{}{
		"id":      c.Response().Header().Get(echo.HeaderXRequestID),
		"message": msg,
		"status":  code,
	}

	c.Logger().Errorj(content)

	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD {
			c.NoContent(code)
		} else {
			c.JSON(code, &strut.ErrResponse{strut.ErrContent{code, msg}})
		}
	}
}
