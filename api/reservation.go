package api

import (
	gen "bitbucket.org/ricardomvpinto/stock-service/api/structures"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

// Validates the consistency of the Reservation struct
func ValidateReservation(res gen.Reservation) error {
	if res.Sku == "" {
		return fmt.Errorf("Sku is empty")
	}
	if res.Warehouse == "" {
		return fmt.Errorf("Warehouse is empty")
	}
	return nil
}

// Processes a Reservation request
func ProcessRequest(r gen.Reservation, put bool, rp gen.RepositoryDefinition, p gen.PubSub) (int, error) {
	var skuFound *gen.Sku

	if err := ValidateReservation(r); err != nil {
		return http.StatusBadRequest, err
	}

	skuFound, err := rp.RepoFindBySkuAndWharehouse(r.Sku, r.Warehouse)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if skuFound.Sku != "" {
		if put {
			if err := rp.RepoInsertReservation(r); err != nil {
				return http.StatusInternalServerError, err
			}
		} else {
			if err := rp.RepoDeleteReservation(r); err != nil {
				if err.Error() == "404" {
					return http.StatusNotFound, fmt.Errorf("No reservation found for Sku %s and Warehouse %s", r.Sku, r.Warehouse)
				}
				return http.StatusInternalServerError, err
			}
		}
	} else {
		return http.StatusNotFound, fmt.Errorf("Sku not found")
	}

	skuResponse, err := rp.RepoFindSku(r.Sku)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("Sku %s not found", r.Sku)
	}

	if err := p.Publish(skuResponse); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// Handler to PUT Reservation request
func PutReservation(rp gen.RepositoryDefinition, p gen.PubSub) echo.HandlerFunc {
	return func(c echo.Context) error {
		var res gen.Reservation

		if err := c.Bind(&res); err != nil {
			return c.JSON(http.StatusBadRequest, &ErrResponse{ErrContent{http.StatusBadRequest, err.Error()}})
		}

		if res.Sku = c.Param("sku"); res.Sku == "" {
			return c.JSON(http.StatusBadRequest, &ErrResponse{ErrContent{http.StatusBadRequest, "Sku not set"}})
		}

		if code, err := ProcessRequest(res, true, rp, p); err != nil {
			return c.JSON(code, &ErrResponse{ErrContent{code, err.Error()}})
		}

		return c.NoContent(http.StatusOK)
	}
}

// Handler to DELETE Reservation request
func RemoveReservation(rp gen.RepositoryDefinition, p gen.PubSub) echo.HandlerFunc {
	return func(c echo.Context) error {
		var res gen.Reservation

		if err := c.Bind(&res); err != nil {
			return c.JSON(http.StatusBadRequest, &ErrResponse{ErrContent{http.StatusBadRequest, err.Error()}})
		}

		if res.Sku = c.Param("sku"); res.Sku == "" {
			return c.JSON(http.StatusBadRequest, &ErrResponse{ErrContent{http.StatusBadRequest, "Sku not set"}})
		}

		if code, err := ProcessRequest(res, false, rp, p); err != nil {
			return c.JSON(code, &ErrResponse{ErrContent{code, err.Error()}})
		}

		return c.NoContent(http.StatusOK)
	}
}
