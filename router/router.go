package router

import (
	gen "bitbucket.org/ricardomvpinto/stock-service/api/structures"
	cnf "bitbucket.org/ricardomvpinto/stock-service/config"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(routes gen.Routes) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = cnf.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
