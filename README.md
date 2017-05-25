# Stock service
Stock service is a small app to deal with stock and stock reservation
The database used to store the data is a mysql one

## Requirements
App requires Golang 1.8 or later, Glide Package Manager and Docker (for building)

## Installation
- Install [Golang](https://golang.org/doc/install)
- Install [Glide](https://glide.sh)
- Install [Docker](htts://docker.com)


## Build
For building binaries please use make, look at the commands bellow:


To manually build the app use:
```
CGO_ENABLED=0 go build -o ./build/stock-service -ldflags "-s -w" -tags netgo -a -v
```

```
// Build the binary in your environment
$ make build

// Build with another OS. Default Linux
$ make OS=darwin build

// Build with custom version.
$ make APP_VERSION=0.1.0 build

// Build with custom app name.
$ make APP_NAME=catalog-search build

// Passing all flags
$ make OS=darwin APP_NAME=stock-service APP_VERSION=0.1.0 build

// Clean Up
$ make clean

// Configure. Install app dependencies.
$ make configure

// Check if docker exists.
$ make depend

// Create a docker image with application
$ make pack

// Pack with custom Docker namespace. Default gfgit
$ make DOCKER_NS=gfgit pack

// Pack with custom version.
$ make APP_VERSION=0.1.0 pack

// Pack with custom app name.
$ make APP_NAME=stock-service pack

// Pack passing all flags
$ make APP_NAME=stock-service APP_VERSION=0.1.0 DOCKER_NS=gfgit pack
```

## Development
```
// Running tests
$ make test

// Running tests with coverage. Output coverage file: coverage.html
$ make test-coverage

// Running tests with junit report. Output coverage file: report.xml
$ make test-report
```

## Usage:

* PUT RESERVATION CALL
curl -v -X PUT http://localhost:8080/reservation/ABCDE -H 'content-type: application/json' -d '{"warehouse":"B"}'

* REMOVE RESERVATION CALL
curl -v -X DELETE http://localhost:8080/reservation/ABCDE -H 'content-type: application/json' -d '{"warehouse":"B"}'

* PUT STOCK CALL
curl -v -X PUT http://localhost:8080/stock/ABCDE -H 'content-type: application/json' -d '{"quantity":20,"warehouse":"B"}'

* GET STOCK CALL
curl -v -X GET http://localhost:8080/stock/ABCDE