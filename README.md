# Stock service
Stock service is a small app to deal with stock and stock reservation
The database used to store the data is a mysql one

## Requirements
App requires Golang 1.9 or later, Glide Package Manager and Docker (for building)

## Installation
- Install [Golang](https://golang.org/doc/install)
- Install [Glide](https://glide.sh)
- Install [Docker](https://www.docker.com/)
- Install [Docker-compose](https://docs.docker.com/compose/)

## Build
For building binaries please use make, look at the commands bellow:

```
// Build the binary in your environment
$ make build

// Build with another OS. Default Linux
$ make OS=darwin build

// Build with custom version.
$ make APP_VERSION=0.1.0 build

// Build with custom app name.
$ make APP_NAME=stock-service build

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

// Pack with custom version.
$ make APP_VERSION=0.1.0 pack

// Pack with custom app name.
$ make APP_NAME=stock-service pack

// Pack passing all flags
$ make APP_NAME=stock-service APP_VERSION=0.1.0
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

## Run it
```
// Build and run docker-compose
$ make build; sudo docker-compose up;

// Run the service
$ ./build/stock-service -l 0.0.0.0:8080 -d core.database.yml.example -p core.rabbitmq.yml.example
```

## Usage:

* PUT RESERVATION CALL
```
curl -v -X PUT http://localhost:8080/reservation/ABCDE -H 'content-type: application/json' -d '{"warehouse":"B"}'
```
* REMOVE RESERVATION CALL
```
curl -v -X DELETE http://localhost:8080/reservation/ABCDE -H 'content-type: application/json' -d '{"warehouse":"B"}'
```
* PUT STOCK CALL
```
curl -v -X PUT http://localhost:8080/stock/ABCDE -H 'content-type: application/json' -d '{"quantity":20,"warehouse":"B"}'
```

* GET STOCK CALL
```
curl -v -X GET http://localhost:8080/stock/ABCDE
```