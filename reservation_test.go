package main

import (
	"net/http"
	"net/http/httptest"
	"fmt"
	"io"
	"bytes"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/gorilla/mux"
)

/* ValidateReservation DataProvider */
type testReserv struct {
	value  Reservation
  	result error
}
var testValidReservation = []testReserv {
	{Reservation{"","AB"}, fmt.Errorf("Sku is empty")},
	{Reservation{"AA",""}, fmt.Errorf("Warehouse is empty")},
	{Reservation{"AA","AB"}, nil},
}

/* Test for ValidateSku method */
func TestValidateReservation(t *testing.T) {
	for _, pair := range testValidReservation {
		v := ValidateReservation(pair.value)
		assert.Equal(t, v, pair.result,"Error message doesn't match")
	}
}

// Mock Server
func setupServerRes() (*mux.Router, *httptest.ResponseRecorder) {
	var routes = Routes{
		Route{
			"PutReservation",
			"PUT",
			"/reservation/{sku}",
			PutReservation(new(RepositoryMock), new(PublisherMock)),
		},Route{
			"RemoveReservation",
			"DELETE",
			"/reservation/{sku}",
			RemoveReservation(new(RepositoryMock), new(PublisherMock)),
		},
	}
    //mux router with added question routes
    m := NewRouter(routes)
    //The response recorder used to record HTTP responses
    respRec := httptest.NewRecorder()

    return m, respRec
}

/* 
Tests for PutReservation method 
*/
type reservationProvider struct {
    method string
    value string
    json string
    result int
}
var testReservationProvider = []reservationProvider {
	{"PUT", "/reservation/", "", http.StatusNotFound}, // url not found
	{"PUT", "/reservation/SAC", `{}`, http.StatusBadRequest}, // invalid Reservation object
	{"PUT", "/reservation/SAC", `{"warehouse":"A"}`, http.StatusInternalServerError}, // RepoFindBySkuAndWharehouse error
	{"PUT", "/reservation/SC", `{"warehouse":"C"}`, http.StatusInternalServerError}, // RepoInsertReservation error
	{"PUT", "/reservation/SC", `{"warehouse":"B"}`, http.StatusNotFound}, // Sku and Warehouse not found
	{"PUT", "/reservation/SC", `{"warehouse":"A"}`, http.StatusOK}, // Insert OK
	{"DELETE", "/reservation/", "", http.StatusNotFound}, // url not found
	{"DELETE", "/reservation/SAC", `{}`, http.StatusBadRequest}, // invalid Reservation object
	{"DELETE", "/reservation/SC", `{"warehouse":"C"}`, http.StatusNotFound}, // Insert OK
	{"DELETE", "/reservation/SC", `{"warehouse":"D"}`, http.StatusInternalServerError}, // Insert OK
	{"DELETE", "/reservation/SC", `{"warehouse":"A"}`, http.StatusOK}, // Insert OK
}

func TestPutDeleteReservation(t *testing.T) {
	for _, pair := range testReservationProvider {
		m, rr := setupServerRes()

		var val io.Reader
		if pair.json != "" {
			var jsonStr = []byte(pair.json)
			val = bytes.NewBuffer(jsonStr)
		}

	    req, err := http.NewRequest(pair.method, pair.value, val)
	    req.Header.Set("Content-Type", "application/json")
	    if err != nil {
	        t.Fatal("TestPutReservation failed!")
	    }

	    m.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, pair.result, "Code doesn't match")
	}
}
