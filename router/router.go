package router

import (
	"github.com/gorilla/mux"
	"net/http"
	gen "bitbucket.org/ricardomvpinto/stock-service/general"
)

func NewRouter(routes gen.Routes) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
