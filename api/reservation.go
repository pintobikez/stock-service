package api

import (
	gen "bitbucket.org/ricardomvpinto/stock-service/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
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
func ProcessRequest(w http.ResponseWriter, r gen.Reservation, put bool, rp gen.RepositoryDefinition, p gen.PubSub) (int, error) {
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
func PutReservation(rp gen.RepositoryDefinition, p gen.PubSub) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var res gen.Reservation
		var isset bool
		_ = json.NewDecoder(r.Body).Decode(&res)

		if res.Sku, isset = vars["sku"]; !isset {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		code, err := ProcessRequest(w, res, true, rp, p)
		w.WriteHeader(code)
		if err != nil {
			json.NewEncoder(w).Encode(gen.JsonErr{Code: code, Text: err.Error()})
		}

		return
	}
}

// Handler to DELETE Reservation request
func RemoveReservation(rp gen.RepositoryDefinition, p gen.PubSub) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var res gen.Reservation
		var isset bool
		_ = json.NewDecoder(r.Body).Decode(&res)

		if res.Sku, isset = vars["sku"]; !isset {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if code, err := ProcessRequest(w, res, false, rp, p); err != nil {
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(gen.JsonErr{Code: code, Text: err.Error()})
		} else {
			w.WriteHeader(code)
		}

		return
	}
}
