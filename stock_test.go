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

/* ValidateSKu DataProvider */
type testSku struct {
	value  Sku
  	result error
}
var testValidSku = []testSku {
	{Sku{"",10,"AB"}, fmt.Errorf("Sku is empty")},
	{Sku{"AA",10,""}, fmt.Errorf("Warehouse is empty")},
	{Sku{"AA",-1,"AB"}, fmt.Errorf("Quantity is negative")},
	{Sku{"AA",10,"AB"}, nil},
}

/* Test for ValidateSku method */
func TestValidateSku(t *testing.T) {
	for _, pair := range testValidSku {
		v := ValidateSku(pair.value)
		assert.Equal(t, v, pair.result,"Error message doesn't match")
	}
}

//Mock Server
func setupServerStock() (*mux.Router, *httptest.ResponseRecorder) {
	var routes = Routes{
		Route{
			"PutStock",
			"PUT",
			"/stock/{sku}",
			PutStock(new(RepositoryMock), new(PublisherMock)),
		},Route{
			"GetStock",
			"GET",
			"/stock/{sku}",
			GetStock(new(RepositoryMock)),
		},
	}
    //mux router with added question routes
    m := NewRouter(routes)
    //The response recorder used to record HTTP responses
    respRec := httptest.NewRecorder()

    return m, respRec
}

/* 
Tests for GetStock method 
*/
type getStockProvider struct {
    value string
    result int
}
var testGetStockProvider = []getStockProvider {
	{"/stock/", http.StatusNotFound}, // url not found
	{"/stock/SAC", http.StatusNotFound}, // sku not found
	{"/stock/SC", http.StatusOK}, // sku found
}
func TestGetStock(t *testing.T) {
	for _, pair := range testGetStockProvider {
		m, rr := setupServerStock()
	    req, err := http.NewRequest("GET", pair.value, nil)
	    if err != nil {
	        t.Fatal("TestGetStock failed!")
	    }

	    m.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, pair.result, "Code doesn't match")
	}
}


/* 
Test for PutStock method
*/
type putStockProvider struct {
    value string
    json string
    result int
}
var testPutStockProvider = []putStockProvider {
	{"/stock/", "", http.StatusNotFound},
	{"/stock/SAC", `{"quantity":10}`, http.StatusBadRequest},
	{"/stock/SAC", `{"quantity":10, "warehouse":"C"}`, http.StatusInternalServerError}, // UPDATE NOK
	{"/stock/SC", `{"quantity":10, "warehouse":"A"}`, http.StatusOK}, // UPDATE OK
	{"/stock/SC", `{"quantity":10, "warehouse":"B"}`, http.StatusOK}, // INSERT OK
	{"/stock/SC", `{"quantity":10, "warehouse":"C"}`, http.StatusInternalServerError}, // UPDATE NOK
	{"/stock/SC", `{"quantity":10, "warehouse":"D"}`, http.StatusInternalServerError}, // INSERT NOK
}
func TestPutStock(t *testing.T) {
	for _, pair := range testPutStockProvider {
		m, rr := setupServerStock()

		var val io.Reader
		if pair.json != "" {
			var jsonStr = []byte(pair.json)
			val = bytes.NewBuffer(jsonStr)
		}

	    req, err := http.NewRequest("PUT", pair.value, val)
	    req.Header.Set("Content-Type", "application/json")
	    if err != nil {
	        t.Fatal("TestPutStock failed!")
	    }

	    m.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, pair.result, "Code doesn't match")
	}
}