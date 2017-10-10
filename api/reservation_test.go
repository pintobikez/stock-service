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

/* ValidateReservation DataProvider */
type testReserv struct {
	value  gen.Reservation
	result error
}

var testValidReservation = []testReserv{
	{gen.Reservation{"", "AB"}, fmt.Errorf("Sku is empty")},
	{gen.Reservation{"AA", ""}, fmt.Errorf("Warehouse is empty")},
	{gen.Reservation{"AA", "AB"}, nil},
}

/* Test for ValidateSku method */
func TestValidateReservation(t *testing.T) {
	for _, pair := range testValidReservation {
		v := ValidateReservation(&pair.value)
		assert.Equal(t, v, pair.result, "Error message doesn't match")
	}
}

/*
Tests for PutReservation method
*/
type reservationProvider struct {
	method string
	value  string
	json   string
	result int
}

var testReservationProvider = []reservationProvider{
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
	for _, pair := range testReservationProvider {
		p := new(mock.PublisherMock)
		r := new(mock.RepositoryMock)

		// Setup
		e := echo.New()
		if pair.method == "PUT" {
			e.PUT("/reservation/:sku", PutReservation(r, p))
		}
		if pair.method == "DELETE" {
			e.DELETE("/reservation/:sku", RemoveReservation(r, p))
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(pair.method, pair.value, strings.NewReader(pair.json))
		req.Header.Set("Content-Type", "application/json")
		e.ServeHTTP(rec, req)

		assert.Equal(t, pair.result, rec.Code, "Code doesn't match")
	}
}
