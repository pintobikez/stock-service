package api

import (
	gen "bitbucket.org/ricardomvpinto/stock-service/api/structures"
	mock "bitbucket.org/ricardomvpinto/stock-service/mocks"
	"fmt"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

/* ValidateSKu DataProvider */
type testSku struct {
	value  gen.Sku
	result error
}

var testValidSku = []testSku{
	{gen.Sku{"", 10, "AB"}, fmt.Errorf("Sku is empty")},
	{gen.Sku{"AA", 10, ""}, fmt.Errorf("Warehouse is empty")},
	{gen.Sku{"AA", -1, "AB"}, fmt.Errorf("Quantity is negative")},
	{gen.Sku{"AA", 10, "AB"}, nil},
}

/* Test for ValidateSku method */
func TestValidateSku(t *testing.T) {
	for _, pair := range testValidSku {
		v := ValidateSku(&pair.value)
		assert.Equal(t, v, pair.result, "Error message doesn't match")
	}
}

/*
Tests for GetStock method
*/
type getStockProvider struct {
	value  string
	result int
}

var testGetStockProvider = []getStockProvider{
	{"/stock/", http.StatusNotFound},    // url not found
	{"/stock/SCA", http.StatusNotFound}, // sku not found
	{"/stock/SC", http.StatusOK},        // sku found
}

func TestGetStock(t *testing.T) {
	for _, pair := range testGetStockProvider {
		r := new(mock.RepositoryMock)

		// Setup
		e := echo.New()
		e.GET("/stock/:sku", GetStock(r))

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", pair.value, strings.NewReader(""))
		req.Header.Set("Content-Type", "application/json")
		e.ServeHTTP(rec, req)
		assert.Equal(t, pair.result, rec.Code, "Code doesn't match")
	}
}

/*
Test for PutStock method
*/
type putStockProvider struct {
	value  string
	json   string
	result int
}

var testPutStockProvider = []putStockProvider{
	{"/stock/", "", http.StatusNotFound},
	{"/stock/SAC", `{"quantity":10}`, http.StatusBadRequest},
	{"/stock/SAC", `{"quantity":10, "warehouse":"C"}`, http.StatusInternalServerError}, // RepoFindBySkuAndWharehouse error
	{"/stock/DDD", `{"quantity":10, "warehouse":"A"}`, http.StatusInternalServerError}, // RepoFindBySkuAndWharehouse Sku empty, INSERT erro
	{"/stock/DDDD", `{"quantity":10, "warehouse":"A"}`, http.StatusOK},                 // RepoFindBySkuAndWharehouse Sku empty, INSERT OK
	{"/stock/SC", `{"quantity":10, "warehouse":"C"}`, http.StatusInternalServerError},  // UPDATE NOK
	{"/stock/SCCC", `{"quantity":10, "warehouse":"B"}`, http.StatusNotFound},           // FindSku to publish error
	{"/stock/SCD", `{"quantity":10, "warehouse":"D"}`, http.StatusInternalServerError}, // Error in publish
}

func TestPutStock(t *testing.T) {
	for _, pair := range testPutStockProvider {
		p := new(mock.PublisherMock)
		r := new(mock.RepositoryMock)

		// Setup
		e := echo.New()
		e.PUT("/stock/:sku", PutStock(r, p))

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("PUT", pair.value, strings.NewReader(pair.json))
		req.Header.Set("Content-Type", "application/json")
		e.ServeHTTP(rec, req)
		assert.Equal(t, pair.result, rec.Code, "Code doesn't match")
	}
}
