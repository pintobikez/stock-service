package main

import (
	"log"
	"strconv"
	"io/ioutil"
	"net/http"
	"gopkg.in/yaml.v2"
)

var repo *Repository = new(Repository)
var pubsub *FilePublisher = new(FilePublisher)

func loadYaml() string {
	var stringConn = ""
	t := Yconfig{}
    
    data, err := ioutil.ReadFile("config/dev.yml")
    if err != nil {
        panic(err)
    }
    err = yaml.Unmarshal([]byte(data), &t)
    if err != nil {
        panic(err)
    }

	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	stringConn = t.Mysql.User + ":" + t.Mysql.Pw
	stringConn += "@tcp(" + t.Mysql.Host + ":" + strconv.Itoa(t.Mysql.Port) +")"
	stringConn += "/" + t.Mysql.Schema + "?charset=utf8"

	return stringConn
}

func main() {
	stringConn := loadYaml()

	var routes = Routes{
		Route{
			"PutStock",
			"PUT",
			"/stock/{sku}",
			PutStock(repo,pubsub),
		},
		Route{
			"GetStock",
			"GET",
			"/stock/{sku}",
			GetStock(repo),
		},
		Route{
			"PutReservation",
			"PUT",
			"/reservation/{sku}",
			PutReservation(repo,pubsub),
		},
		Route{
			"RemoveReservation",
			"DELETE",
			"/reservation/{sku}",
			RemoveReservation(repo,pubsub),
		},
	}

	repo.connectDB(stringConn)
	router := NewRouter(routes)
	log.Fatal(http.ListenAndServe(":8080", router))

	defer repo.disconnectDB()
}
