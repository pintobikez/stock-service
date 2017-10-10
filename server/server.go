package server

import (
	"github.com/coreos/go-systemd/activation"
	"github.com/labstack/echo"
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
