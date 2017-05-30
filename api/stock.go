package api

import (
	gen "bitbucket.org/ricardomvpinto/stock-service/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
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
func GetStock(rp gen.RepositoryDefinition) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		vars := mux.Vars(r)
		var skuResponse *gen.SkuResponse
		var err error

		skuValue, isset := vars["sku"]
		if !isset {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		skuResponse, err = rp.RepoFindSku(skuValue)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			if err2 := json.NewEncoder(w).Encode(gen.JsonErr{Code: http.StatusNotFound, Text: err.Error()}); err2 != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		if err = json.NewEncoder(w).Encode(skuResponse); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// Handler to PUT Stock request
func PutStock(rp gen.RepositoryDefinition, p gen.PubSub) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		var err error
		var isset bool
		var af int64 = 1

		vars := mux.Vars(r)
		var s gen.Sku
		_ = json.NewDecoder(r.Body).Decode(&s)

		if s.Sku, isset = vars["sku"]; !isset {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := ValidateSku(s); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(gen.JsonErr{Code: http.StatusBadRequest, Text: err.Error()})
			return
		}

		f, erre := rp.RepoFindBySkuAndWharehouse(s.Sku, s.Warehouse)
		if erre != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(gen.JsonErr{Code: http.StatusInternalServerError, Text: erre.Error()})
			return
		}

		if f.Sku != "" {
			if af, err = rp.RepoUpdateSku(s); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(gen.JsonErr{Code: http.StatusInternalServerError, Text: err.Error()})
				return
			}
		} else {
			if err = rp.RepoInsertSku(s); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(gen.JsonErr{Code: http.StatusInternalServerError, Text: err.Error()})
				return
			}
		}

		json.NewEncoder(w).Encode(s)
		if af > 0 { //publish message
			skuResponse, err := rp.RepoFindSku(s.Sku)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(gen.JsonErr{Code: http.StatusNotFound, Text: "Sku " + s.Sku + " not found"})
				return
			}

			if err := p.Publish(skuResponse); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(gen.JsonErr{Code: http.StatusInternalServerError, Text: err.Error()})
				return
			}
		}
		return
	}
}
