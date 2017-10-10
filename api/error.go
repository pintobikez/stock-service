package api

import (
	"github.com/labstack/echo"
	"net/http"
)

var (
	SkuNotFound            = "Sku %s not found"
	ReservationDeleteError = "No reservation found for Sku %s and Warehouse %s"
)

type (
	ErrResponse struct {
		Error ErrContent `json:"error"`
	}

	ErrContent struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
)

func Error(err error, c echo.Context) {
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
			c.JSON(code, &ErrResponse{ErrContent{code, msg}})
		}
	}
}
