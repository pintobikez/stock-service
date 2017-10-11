package api

import (
	gen "bitbucket.org/ricardomvpinto/stock-service/api/structures"
	pub "bitbucket.org/ricardomvpinto/stock-service/publisher"
	repo "bitbucket.org/ricardomvpinto/stock-service/repository"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

type Api struct {
	rp repo.IRepository
	pb pub.IPubSub
}

func New(rpo repo.IRepository, p pub.IPubSub) *Api {
	return &Api{rp: rpo, pb: p}
}

// Handler to GET Stock request
func (a *Api) GetStock() echo.HandlerFunc {
	return func(c echo.Context) error {

		skuValue := c.Param("sku")
		skuResponse, err := a.rp.RepoFindSku(skuValue)

		if err != nil {
			return c.JSON(http.StatusNotFound, &ErrResponse{ErrContent{http.StatusNotFound, err.Error()}})
		}

		return c.JSON(http.StatusOK, skuResponse)
	}
}

// Handler to PUT Stock request
func (a *Api) PutStock() echo.HandlerFunc {
	return func(c echo.Context) error {

		var af int64 = 1
		var s *gen.Sku

		if err := c.Bind(&s); err != nil {
			return c.JSON(http.StatusBadRequest, &ErrResponse{ErrContent{http.StatusBadRequest, err.Error()}})
		}

		s.Sku = c.Param("sku")

		if err := a.validateSku(s); err != nil {
			return c.JSON(http.StatusBadRequest, &ErrResponse{ErrContent{http.StatusBadRequest, err.Error()}})
		}

		f, err := a.rp.RepoFindBySkuAndWharehouse(s.Sku, s.Warehouse)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &ErrResponse{ErrContent{http.StatusInternalServerError, err.Error()}})
		}

		if f.Sku != "" {
			af, err = a.rp.RepoUpdateSku(s)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, &ErrResponse{ErrContent{http.StatusInternalServerError, err.Error()}})
			}
		} else {
			if err := a.rp.RepoInsertSku(s); err != nil {
				return c.JSON(http.StatusInternalServerError, &ErrResponse{ErrContent{http.StatusInternalServerError, err.Error()}})
			}
		}

		if af > 0 { //publish message
			skuResponse, err := a.rp.RepoFindSku(s.Sku)
			if err != nil {
				return c.JSON(http.StatusNotFound, &ErrResponse{ErrContent{http.StatusInternalServerError, fmt.Sprintf(SkuNotFound, s.Sku)}})
			}

			if err := a.pb.Publish(skuResponse); err != nil {
				return c.JSON(http.StatusInternalServerError, &ErrResponse{ErrContent{http.StatusInternalServerError, err.Error()}})
			}
		}

		return c.NoContent(http.StatusOK)
	}
}

// Handler to PUT Reservation request
func (a *Api) PutReservation() echo.HandlerFunc {
	return func(c echo.Context) error {
		var res *gen.Reservation

		if err := c.Bind(&res); err != nil {
			return c.JSON(http.StatusBadRequest, &ErrResponse{ErrContent{http.StatusBadRequest, err.Error()}})
		}
		res.Sku = c.Param("sku")

		if code, err := a.processRequest(res, true); err != nil {
			return c.JSON(code, &ErrResponse{ErrContent{code, err.Error()}})
		}

		return c.NoContent(http.StatusOK)
	}
}

// Handler to DELETE Reservation request
func (a *Api) RemoveReservation() echo.HandlerFunc {
	return func(c echo.Context) error {
		var res *gen.Reservation

		if err := c.Bind(&res); err != nil {
			return c.JSON(http.StatusBadRequest, &ErrResponse{ErrContent{http.StatusBadRequest, err.Error()}})
		}
		res.Sku = c.Param("sku")

		if code, err := a.processRequest(res, false); err != nil {
			return c.JSON(code, &ErrResponse{ErrContent{code, err.Error()}})
		}

		return c.NoContent(http.StatusOK)
	}
}

// Processes a Reservation request
func (a *Api) processRequest(r *gen.Reservation, put bool) (int, error) {
	var skuFound *gen.Sku

	if err := a.validateReservation(r); err != nil {
		return http.StatusBadRequest, err
	}

	skuFound, err := a.rp.RepoFindBySkuAndWharehouse(r.Sku, r.Warehouse)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if skuFound.Sku != "" {
		if put {
			if err := a.rp.RepoInsertReservation(r); err != nil {
				return http.StatusInternalServerError, err
			}
		} else {
			if err := a.rp.RepoDeleteReservation(r); err != nil {
				if err.Error() == "404" {
					return http.StatusNotFound, fmt.Errorf(ReservationDeleteError, r.Sku, r.Warehouse)
				}
				return http.StatusInternalServerError, err
			}
		}
	} else {
		return http.StatusNotFound, fmt.Errorf(SkuNotFound, "")
	}

	skuResponse, err := a.rp.RepoFindSku(r.Sku)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf(SkuNotFound, r.Sku)
	}

	if err := a.pb.Publish(skuResponse); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// Validates the consistency of the Sku struct
func (a *Api) validateSku(s *gen.Sku) error {
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
func (a *Api) validateReservation(res *gen.Reservation) error {
	if res.Sku == "" {
		return fmt.Errorf("Sku is empty")
	}
	if res.Warehouse == "" {
		return fmt.Errorf("Warehouse is empty")
	}
	return nil
}
