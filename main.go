package main

import (
	"os"
	"log"
	"strconv"
	"io/ioutil"
	"net/http"
	"gopkg.in/yaml.v2"
        gen "bitbucket.org/ricardomvpinto/stock-service/utils"
        rou "bitbucket.org/ricardomvpinto/stock-service/router"
	api "bitbucket.org/ricardomvpinto/stock-service/api"	
	pub "bitbucket.org/ricardomvpinto/stock-service/publisher"
	rep "bitbucket.org/ricardomvpinto/stock-service/repository"
)

var repo *rep.Repository = new(rep.Repository)
var pubsub *pub.FilePublisher = new(pub.FilePublisher)

func buildStringConnection(filename string) string {
	t, err := gen.loadConfigFile(filename)
	if err != nill {
		panic(err)
	}
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	stringConn = t.Driver.User + ":" + t.Driver.Pw
	stringConn += "@tcp(" + t.Driver.Host + ":" + strconv.Itoa(t.Driver.Port) +")"
	stringConn += "/" + t.Driver.Schema + "?charset=utf8"

	return stringConn
}

func main() {
	stringConn := buildStringConnection(os.Getenv(STORAGE_FILE))

	var routes = gen.Routes{
		gen.Route{
			"PutStock",
			"PUT",
			"/stock/{sku}",
			api.PutStock(repo,pubsub),
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
			api.PutReservation(repo,pubsub),
		},
		gen.Route{
			"RemoveReservation",
			"DELETE",
			"/reservation/{sku}",
			api.RemoveReservation(repo,pubsub),
		},
	}

	repo.connectDB(stringConn)
	router := rou.NewRouter(routes)
	log.Fatal(http.ListenAndServe(":8080", router))

	defer repo.disconnectDB()
}
