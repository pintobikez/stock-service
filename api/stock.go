package api

import (
	gen "bitbucket.org/ricardomvpinto/stock-service/api/structures"
	//"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

// Validates the consistency of the Sku struct
func ValidateSku(s gen.Sku) error {
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
func GetStock(rp gen.RepositoryDefinition) echo.HandlerFunc {
	return func(c echo.Context) error {
		var skuResponse *gen.SkuResponse
		var err error

		skuValue := c.Param("sku")
		if skuValue == "" {
			return c.JSON(http.StatusBadRequest, "Sku not set")
		}

		skuResponse, err = rp.RepoFindSku(skuValue)
		if err != nil {
			return c.JSON(http.StatusNotFound, gen.JsonErr{Code: http.StatusNotFound, Text: err.Error()})
		}

		/*if res, err2 := json.Marshal(skuResponse); err2 != nil {
			return c.JSON(http.StatusInternalServerError, gen.JsonErr{Code: http.StatusInternalServerError, Text: err2.Error()})
		}*/

		return c.JSON(http.StatusOK, skuResponse)
	}
}

// Handler to PUT Stock request
func PutStock(rp gen.RepositoryDefinition, p gen.PubSub) echo.HandlerFunc {
	return func(c echo.Context) error {

		var err error
		var af int64 = 1

		var s gen.Sku
		if err := c.Bind(&s); err != nil {
			return c.JSON(http.StatusBadRequest, gen.JsonErr{Code: http.StatusBadRequest, Text: err.Error()})
		}

		if s.Sku = c.Param("sku"); s.Sku == "" {
			return c.JSON(http.StatusBadRequest, "Sku not set")
		}

		if err := ValidateSku(s); err != nil {
			return c.JSON(http.StatusBadRequest, gen.JsonErr{Code: http.StatusBadRequest, Text: err.Error()})
		}

		f, erre := rp.RepoFindBySkuAndWharehouse(s.Sku, s.Warehouse)
		if erre != nil {
			return c.JSON(http.StatusInternalServerError, gen.JsonErr{Code: http.StatusInternalServerError, Text: err.Error()})
		}

		if f.Sku != "" {
			if af, err = rp.RepoUpdateSku(s); err != nil {
				return c.JSON(http.StatusInternalServerError, gen.JsonErr{Code: http.StatusInternalServerError, Text: err.Error()})
			}
		} else {
			if err = rp.RepoInsertSku(s); err != nil {
				return c.JSON(http.StatusInternalServerError, gen.JsonErr{Code: http.StatusInternalServerError, Text: err.Error()})
			}
		}

		if af > 0 { //publish message
			skuResponse, err := rp.RepoFindSku(s.Sku)
			if err != nil {
				return c.JSON(http.StatusNotFound, gen.JsonErr{Code: http.StatusNotFound, Text: "Sku " + s.Sku + " not found"})
			}

			if err := p.Publish(skuResponse); err != nil {
				return c.JSON(http.StatusInternalServerError, gen.JsonErr{Code: http.StatusInternalServerError, Text: err.Error()})
			}
		}

		return c.NoContent(http.StatusOK)
	}
}
