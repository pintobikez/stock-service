package api

import (
	strut "bitbucket.org/ricardomvpinto/stock-service/api/structures"
	pub "bitbucket.org/ricardomvpinto/stock-service/publisher"
	repo "bitbucket.org/ricardomvpinto/stock-service/repository"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

const (
	StatusAvailable   = "Available"
	StatusUnavailable = "Unavailable"
)

type API struct {
	rp repo.IRepository
	pb pub.IPubSub
}

func New(rpo repo.IRepository, p pub.IPubSub) *API {
	return &API{rp: rpo, pb: p}
}

// Handler for Health Status
func (a *API) HealthStatus() echo.HandlerFunc {
	return func(c echo.Context) error {

		resp := &strut.HealthStatus{
			Pub:  &strut.HealthStatusDetail{Status: StatusAvailable, Detail: ""},
			Repo: &strut.HealthStatusDetail{Status: StatusAvailable, Detail: ""},
		}

		if err := a.pb.Health(); err != nil {
			resp.Pub.Status = StatusUnavailable
			resp.Pub.Detail = err.Error()
		}
		if err := a.rp.Health(); err != nil {
			resp.Repo.Status = StatusUnavailable
			resp.Repo.Detail = err.Error()
		}

		return c.JSON(http.StatusOK, resp)
	}
}

// Handler to GET Stock request
func (a *API) GetStock() echo.HandlerFunc {
	return func(c echo.Context) error {

		skuValue := c.Param("sku")
		skuResponse, err := a.rp.RepoFindSku(skuValue)

		if err != nil {
			return c.JSON(http.StatusNotFound, &ErrResponse{ErrContent{ErrorCodeSkuNotFound, err.Error()}})
		}

		return c.JSON(http.StatusOK, skuResponse)
	}
}

// Handler to PUT Stock request
func (a *API) PutStock() echo.HandlerFunc {
	return func(c echo.Context) error {

		var af int64 = 1
		var s *strut.Sku

		if err := c.Bind(&s); err != nil {
			return c.JSON(http.StatusBadRequest, &ErrResponse{ErrContent{ErrorCodeWrongJsonFormat, err.Error()}})
		}

		s.Sku = c.Param("sku")

		if err := a.validateSku(s); err != nil {
			return c.JSON(http.StatusBadRequest, &ErrResponse{ErrContent{ErrorCodeInvalidContent, err.Error()}})
		}

		f, err := a.rp.RepoFindBySkuAndWharehouse(s.Sku, s.Warehouse)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &ErrResponse{ErrContent{ErrorCodeSkuNotFound, err.Error()}})
		}

		if f.Sku != "" {
			af, err = a.rp.RepoUpdateSku(s)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, &ErrResponse{ErrContent{ErrorCodeStoringContent, err.Error()}})
			}
		} else {
			if err := a.rp.RepoInsertSku(s); err != nil {
				return c.JSON(http.StatusInternalServerError, &ErrResponse{ErrContent{ErrorCodeStoringContent, err.Error()}})
			}
		}

		if af > 0 { //publish message
			skuResponse, err := a.rp.RepoFindSku(s.Sku)
			if err != nil {
				return c.JSON(http.StatusNotFound, &ErrResponse{ErrContent{ErrorCodeSkuNotFound, fmt.Sprintf(SkuNotFound, s.Sku)}})
			}

			if err := a.pb.Publish(skuResponse); err != nil {
				return c.JSON(http.StatusInternalServerError, &ErrResponse{ErrContent{ErrorCodePublishingMessage, err.Error()}})
			}
		}

		return c.NoContent(http.StatusOK)
	}
}

// Handler to PUT Reservation request
func (a *API) PutReservation() echo.HandlerFunc {
	return func(c echo.Context) error {
		var res *strut.Reservation

		if err := c.Bind(&res); err != nil {
			return c.JSON(http.StatusBadRequest, &ErrResponse{ErrContent{ErrorCodeWrongJsonFormat, err.Error()}})
		}
		res.Sku = c.Param("sku")

		if httpcode, code, err := a.processReservation(res, true); err != nil {
			return c.JSON(httpcode, &ErrResponse{ErrContent{code, err.Error()}})
		}

		return c.NoContent(http.StatusOK)
	}
}

// Handler to DELETE Reservation request
func (a *API) RemoveReservation() echo.HandlerFunc {
	return func(c echo.Context) error {
		var res *strut.Reservation

		if err := c.Bind(&res); err != nil {
			return c.JSON(http.StatusBadRequest, &ErrResponse{ErrContent{ErrorCodeWrongJsonFormat, err.Error()}})
		}
		res.Sku = c.Param("sku")

		if httpcode, code, err := a.processReservation(res, false); err != nil {
			return c.JSON(httpcode, &ErrResponse{ErrContent{code, err.Error()}})
		}

		return c.NoContent(http.StatusOK)
	}
}

// Processes a Reservation request
func (a *API) processReservation(r *strut.Reservation, put bool) (int, int, error) {
	var skuFound *strut.Sku

	if err := a.validateReservation(r); err != nil {
		return http.StatusBadRequest, ErrorCodeInvalidContent, err
	}

	skuFound, err := a.rp.RepoFindBySkuAndWharehouse(r.Sku, r.Warehouse)
	if err != nil {
		return http.StatusInternalServerError, ErrorCodeSkuNotFound, err
	}

	if skuFound.Sku != "" {
		if put {
			if err := a.rp.RepoInsertReservation(r); err != nil {
				return http.StatusInternalServerError, ErrorCodeStoringContent, err
			}
		} else {
			if err := a.rp.RepoDeleteReservation(r); err != nil {
				if err.Error() == "404" {
					return http.StatusNotFound, ErrorCodeSkuNotFound, fmt.Errorf(ReservationDeleteError, r.Sku, r.Warehouse)
				}
				return http.StatusInternalServerError, ErrorCodeStoringContent, err
			}
		}
	} else {
		return http.StatusNotFound, ErrorCodeSkuNotFound, fmt.Errorf(SkuNotFound, "")
	}

	skuResponse, err := a.rp.RepoFindSku(r.Sku)
	if err != nil {
		return http.StatusNotFound, ErrorCodeSkuNotFound, fmt.Errorf(SkuNotFound, r.Sku)
	}

	if err := a.pb.Publish(skuResponse); err != nil {
		return http.StatusInternalServerError, ErrorCodePublishingMessage, err
	}

	return http.StatusOK, 0, nil
}

// Validates the consistency of the Sku struct
func (a *API) validateSku(s *strut.Sku) error {
	if s.Sku == "" {
		return fmt.Errorf("Sku is empty")
	}
	if s.Warehouse == "" {
		return fmt.Errorf("Warehouse is empty")
	}
	if s.Quantity < 0 {
		return fmt.Errorf("Quantity is negative")
	}
	return nil
}

// Validates the consistency of the Reservation struct
func (a *API) validateReservation(res *strut.Reservation) error {
	if res.Sku == "" {
		return fmt.Errorf("Sku is empty")
	}
	if res.Warehouse == "" {
		return fmt.Errorf("Warehouse is empty")
	}
	return nil
}
