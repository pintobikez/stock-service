package api

import (
	gen "bitbucket.org/ricardomvpinto/stock-service/api/structures"
	pub "bitbucket.org/ricardomvpinto/stock-service/publisher"
	repo "bitbucket.org/ricardomvpinto/stock-service/repository"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

// Validates the consistency of the Sku struct
func ValidateSku(s *gen.Sku) error {
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

// Handler to GET Stock request
func GetStock(rp repo.IRepository) echo.HandlerFunc {
	return func(c echo.Context) error {

		skuValue := c.Param("sku")
		skuResponse, err := rp.RepoFindSku(skuValue)

		if err != nil {
			return c.JSON(http.StatusNotFound, &ErrResponse{ErrContent{http.StatusNotFound, err.Error()}})
		}

		return c.JSON(http.StatusOK, skuResponse)
	}
}

// Handler to PUT Stock request
func PutStock(rp repo.IRepository, p pub.IPubSub) echo.HandlerFunc {
	return func(c echo.Context) error {

		var af int64 = 1
		var s *gen.Sku

		if err := c.Bind(&s); err != nil {
			return c.JSON(http.StatusBadRequest, &ErrResponse{ErrContent{http.StatusBadRequest, err.Error()}})
		}

		s.Sku = c.Param("sku")

		if err := ValidateSku(s); err != nil {
			return c.JSON(http.StatusBadRequest, &ErrResponse{ErrContent{http.StatusBadRequest, err.Error()}})
		}

		f, err := rp.RepoFindBySkuAndWharehouse(s.Sku, s.Warehouse)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &ErrResponse{ErrContent{http.StatusInternalServerError, err.Error()}})
		}

		if f.Sku != "" {
			af, err = rp.RepoUpdateSku(s)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, &ErrResponse{ErrContent{http.StatusInternalServerError, err.Error()}})
			}
		} else {
			if err := rp.RepoInsertSku(s); err != nil {
				return c.JSON(http.StatusInternalServerError, &ErrResponse{ErrContent{http.StatusInternalServerError, err.Error()}})
			}
		}

		if af > 0 { //publish message
			skuResponse, err := rp.RepoFindSku(s.Sku)
			if err != nil {
				return c.JSON(http.StatusNotFound, &ErrResponse{ErrContent{http.StatusInternalServerError, fmt.Sprintf(SkuNotFound, s.Sku)}})
			}

			if err := p.Publish(skuResponse); err != nil {
				return c.JSON(http.StatusInternalServerError, &ErrResponse{ErrContent{http.StatusInternalServerError, err.Error()}})
			}
		}

		return c.NoContent(http.StatusOK)
	}
}
