package main

import "net/http"

type PubSub interface {
	Publish(s *SkuResponse) error
}

type RepositoryDefinition interface {
	connectDB(stringConn string)
	disconnectDB()
	RepoFindBySkuAndWharehouse(sku string, warehouse string) (*Sku, error)
	RepoFindSku(sku string) (*SkuResponse, error)
	RepoUpdateSku(s Sku) (int64, error)
	RepoInsertSku(s Sku) error
	RepoInsertReservation(re Reservation) error
	RepoDeleteReservation(re Reservation) error
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

type Yconfig struct {
    Mysql struct {
        Host string
	    User string
	    Pw string
	    Port int
	    Schema string
    }
}

type jsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

type Sku struct {
	Sku       string `json:"sku"`
	Quantity  int64  `json:"quantity"`
	Warehouse string `json:"warehouse"`
}

type SkuResponse struct {
	Sku       string      `json:"sku"`
	Values    []SkuValues `json:"values"`
	Reserved  int64       `json:"reserved"`
	Available int64       `json:"avail"`
}

type SkuValues struct {
	Quantity  int64  `json:"quantity"`
	Warehouse string `json:"warehouse"`
}

type Reservation struct {
	Sku       string `json:"sku"`
	Warehouse string `json:"warehouse"`
}
