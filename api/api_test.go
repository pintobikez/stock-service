package api

import (
	gen "bitbucket.org/ricardomvpinto/stock-service/api/structures"
	mock "bitbucket.org/ricardomvpinto/stock-service/mocks"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

/*
Tests for PutReservation method
*/
type reservationProviderApi struct {
	method string
	value  string
	json   string
	result int
}

var testReservationProviderApi = []reservationProviderApi{
	{"PUT", "/reservation/", "", http.StatusNotFound},                                   // url not found
	{"PUT", "/reservation/SAC", `{}`, http.StatusBadRequest},                            // invalid Reservation object
	{"PUT", "/reservation/SAC", `{"warehouse":"A"}`, http.StatusInternalServerError},    // RepoFindBySkuAndWharehouse error
	{"PUT", "/reservation/SC", `{"warehouse":"C"}`, http.StatusInternalServerError},     // RepoInsertReservation error
	{"PUT", "/reservation/SCA", `{"warehouse":"B"}`, http.StatusNotFound},               // Sku and Warehouse not found
	{"PUT", "/reservation/SCD", `{"warehouse":"A"}`, http.StatusInternalServerError},    // Error Publish
	{"PUT", "/reservation/SCC", `{"warehouse":"A"}`, http.StatusOK},                     // Insert OK
	{"DELETE", "/reservation/", "", http.StatusNotFound},                                // url not found
	{"DELETE", "/reservation/SAC", `{}`, http.StatusBadRequest},                         // invalid Reservation object
	{"DELETE", "/reservation/SAC", `{"warehouse":"C"}`, http.StatusInternalServerError}, // RepoDeleteReservation error
	{"DELETE", "/reservation/SC", `{"warehouse":"C"}`, http.StatusInternalServerError},  // RepoDeleteReservation error 404
	{"DELETE", "/reservation/SCE", `{"warehouse":"D"}`, http.StatusNotFound},            // RepoDeleteReservation error
	{"DELETE", "/reservation/DDD", `{"warehouse":"D"}`, http.StatusNotFound},            // Sku and Warehouse not found
	{"DELETE", "/reservation/SCC", `{"warehouse":"A"}`, http.StatusOK},                  // Insert OK
}

func TestPutDeleteReservation(t *testing.T) {
	for _, pair := range testReservationProviderApi {
		p := new(mock.PublisherMock)
		r := new(mock.RepositoryMock)
		a := New(r, p)

		// Setup
		e := echo.New()
		if pair.method == "PUT" {
			e.PUT("/reservation/:sku", a.PutReservation())
		}
		if pair.method == "DELETE" {
			e.DELETE("/reservation/:sku", a.RemoveReservation())
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(pair.method, pair.value, strings.NewReader(pair.json))
		req.Header.Set("Content-Type", "application/json")
		e.ServeHTTP(rec, req)

		assert.Equal(t, pair.result, rec.Code, "Code doesn't match")
	}
}

/*
Tests for GetStock method
*/
type getStockProviderApi struct {
	value  string
	result int
}

var testGetStockProviderApi = []getStockProviderApi{
	{"/stock/", http.StatusNotFound},    // url not found
	{"/stock/SCA", http.StatusNotFound}, // sku not found
	{"/stock/SC", http.StatusOK},        // sku found
}

func TestGetStock(t *testing.T) {
	for _, pair := range testGetStockProviderApi {
		p := new(mock.PublisherMock)
		r := new(mock.RepositoryMock)
		a := New(r, p)

		// Setup
		e := echo.New()
		e.GET("/stock/:sku", a.GetStock())

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
type putStockProviderApi struct {
	value  string
	json   string
	result int
}

var testPutStockProviderApi = []putStockProviderApi{
	{"/stock/", "", http.StatusNotFound},                                               // Incorrect url no sku
	{"/stock/SAC", `{"quantity":10}`, http.StatusBadRequest},                           // empty warehouse error
	{"/stock/SAC", `{"quantity":10, "warehouse":"C"}`, http.StatusInternalServerError}, // RepoFindBySkuAndWharehouse error
	{"/stock/DDD", `{"quantity":10, "warehouse":"A"}`, http.StatusInternalServerError}, // RepoFindBySkuAndWharehouse Sku empty, INSERT erro
	{"/stock/DDDD", `{"quantity":10, "warehouse":"A"}`, http.StatusOK},                 // RepoFindBySkuAndWharehouse Sku empty, INSERT OK
	{"/stock/SC", `{"quantity":10, "warehouse":"C"}`, http.StatusInternalServerError},  // UPDATE NOK
	{"/stock/SCCC", `{"quantity":10, "warehouse":"B"}`, http.StatusNotFound},           // FindSku to publish error
	{"/stock/SCD", `{"quantity":10, "warehouse":"D"}`, http.StatusInternalServerError}, // Error in publish
}

func TestPutStock(t *testing.T) {
	for _, pair := range testPutStockProviderApi {
		p := new(mock.PublisherMock)
		r := new(mock.RepositoryMock)
		a := New(r, p)

		// Setup
		e := echo.New()
		e.PUT("/stock/:sku", a.PutStock())

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("PUT", pair.value, strings.NewReader(pair.json))
		req.Header.Set("Content-Type", "application/json")
		e.ServeHTTP(rec, req)
		assert.Equal(t, pair.result, rec.Code, "Code doesn't match")
	}
}

/* ValidateSKu DataProvider */
type testSkuApi struct {
	value  gen.Sku
	result error
}

var testValidSkuApi = []testSkuApi{
	{gen.Sku{"", 10, "AB"}, fmt.Errorf("Sku is empty")},
	{gen.Sku{"AA", 10, ""}, fmt.Errorf("Warehouse is empty")},
	{gen.Sku{"AA", -1, "AB"}, fmt.Errorf("Quantity is negative")},
	{gen.Sku{"AA", 10, "AB"}, nil},
}

/* Test for ValidateSku method */
func TestValidateSku(t *testing.T) {
	p := new(mock.PublisherMock)
	r := new(mock.RepositoryMock)
	a := New(r, p)

	for _, pair := range testValidSkuApi {
		v := a.validateSku(&pair.value)
		assert.Equal(t, v, pair.result, "Error message doesn't match")
	}
}

/* ValidateReservation DataProvider */
type testReservApi struct {
	value  gen.Reservation
	result error
}

var testValidReservationApi = []testReservApi{
	{gen.Reservation{"", "AB"}, fmt.Errorf("Sku is empty")},
	{gen.Reservation{"AA", ""}, fmt.Errorf("Warehouse is empty")},
	{gen.Reservation{"AA", "AB"}, nil},
}

/* Test for ValidateSku method */
func TestValidateReservation(t *testing.T) {
	p := new(mock.PublisherMock)
	r := new(mock.RepositoryMock)
	a := New(r, p)

	for _, pair := range testValidReservationApi {
		v := a.validateReservation(&pair.value)
		assert.Equal(t, v, pair.result, "Error message doesn't match")
	}
}

/*
Tests for HealthStatus method
*/
type getHealthStatusApi struct {
	value string
	erro  string
}

var testGetHealthStatusApi = []getHealthStatusApi{
	{"/health/", "repo"}, // error in repo
	{"/health/", "pub"},  // error in publisher
	{"/health/", ""},     // all good
}

func TestHealthStatus(t *testing.T) {
	for _, pair := range testGetHealthStatusApi {
		p := new(mock.PublisherMock)
		r := new(mock.RepositoryMock)
		a := New(r, p)

		switch pair.erro {
		case "repo":
			r.Iserror = true
			break
		case "pub":
			p.Iserror = true
			break
		}

		// Setup
		e := echo.New()
		e.GET("/health/", a.HealthStatus())

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", pair.value, strings.NewReader(""))
		req.Header.Set("Content-Type", "application/json")
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		val := new(gen.HealthStatus)
		_ = json.Unmarshal([]byte(rec.Body.String()), val)

		// Assertions
		switch pair.erro {
		case "repo":
			assert.Equal(t, StatusUnavailable, val.Repo.Status)
			break
		case "pub":
			assert.Equal(t, StatusUnavailable, val.Pub.Status)
			break
		}
	}
}
