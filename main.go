package main

import (
	api "bitbucket.org/ricardomvpinto/stock-service/api"
	pub "bitbucket.org/ricardomvpinto/stock-service/publisher"
	rep "bitbucket.org/ricardomvpinto/stock-service/repository"
	rou "bitbucket.org/ricardomvpinto/stock-service/router"
	gen "bitbucket.org/ricardomvpinto/stock-service/utils"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	SS_DATABASE_FILE = "SS_DATABASE_FILE"
	SS_LISTEN        = "SS_LISTEN"
)

var repo *rep.Repository = new(rep.Repository)
var pubsub *pub.FilePublisher = new(pub.FilePublisher)

func buildStringConnection(filename string) string {
	t, err := gen.LoadConfigFile(filename)
	if err != nil {
		panic(err)
	}
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	stringConn := t.Driver.User + ":" + t.Driver.Pw
	stringConn += "@tcp(" + t.Driver.Host + ":" + strconv.Itoa(t.Driver.Port) + ")"
	stringConn += "/" + t.Driver.Schema + "?charset=utf8"

	return stringConn
}

func main() {
	stringConn := buildStringConnection(os.Getenv(SS_DATABASE_FILE))

	var routes = gen.Routes{
		gen.Route{
			"PutStock",
			"PUT",
			"/stock/{sku}",
			api.PutStock(repo, pubsub),
		},
		gen.Route{
			"GetStock",
			"GET",
			"/stock/{sku}",
			api.GetStock(repo),
		},
		gen.Route{
			"PutReservation",
			"PUT",
			"/reservation/{sku}",
			api.PutReservation(repo, pubsub),
		},
		gen.Route{
			"RemoveReservation",
			"DELETE",
			"/reservation/{sku}",
			api.RemoveReservation(repo, pubsub),
		},
	}

	repo.ConnectDB(stringConn)
	router := rou.NewRouter(routes)
	log.Fatal(http.ListenAndServe(os.Getenv(SS_LISTEN), router))

	defer repo.DisconnectDB()
}
