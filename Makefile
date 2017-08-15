.PHONY: build clean configure depend pack test test-coverage test-report

APP_NAME=stock-service
APP_PATH=$(shell head -n 1 ./glide.yaml | awk '{print $$2}')
APP_VERSION=0.0.1

LDFLAGS=--ldflags '-X main.version=${APP_VERSION} -X main.appName=${APP_NAME} -extldflags "-static" -w'
OS=linux

DOCKER_NS=gfgit
DOCKER_IMAGE=gfgit/golang-ci:1.8.3

.DEFAULT_GOAL := build

build: depend configure
ifeq (${RUN},docker)
	@docker run --rm \
        -v "$(shell pwd)":/go/src/${APP_PATH} \
        -w /go/src/${APP_PATH} \
        ${DOCKER_IMAGE} sh -c "make OS=${OS} APP_NAME=${APP_NAME} APP_VERSION=${APP_VERSION}"
else
	@CGO_ENABLED=0 GOOS=${OS} go build -a ${LDFLAGS} -tags netgo -installsuffix netgo -v \
    -o ./build/${APP_NAME} bitbucket.org/ricardomvpinto/stock-service/cmd
endif

clean:
ifeq (${RUN},docker)
	@docker run --rm \
        -v "$(shell pwd)":/go/src/${APP_PATH} \
        -w /go/src/${APP_PATH} \
        ${DOCKER_IMAGE} sh -c "make clean"
else
	@rm -fR vendor/ ./build ./.glide/
endif

configure:
ifeq (${RUN},docker)
	@docker run --rm \
        -v "$(shell pwd)":/go/src/${APP_PATH} \
        -w /go/src/${APP_PATH} \
        ${DOCKER_IMAGE} sh -c "make configure"
else
	@glide install
endif

depend:
ifeq (${RUN},docker)
	@docker run --rm \
        -v "$(shell pwd)":/go/src/${APP_PATH} \
        -w /go/src/${APP_PATH} \
        ${DOCKER_IMAGE} sh -c "make depend"
else
	@mkdir -p ./build
endif

pack: depend
	@command -v docker > /dev/null 2>&1 || ( echo "Please install Docker https://docs.docker.com/engine/installation/" && exit 1 )
	@docker build -t ${DOCKER_NS}/${APP_NAME}:${APP_VERSION} --build-arg APP_NAME=${APP_NAME} -f ./Dockerfile .

test:
ifeq (${RUN},docker)
	@docker run --rm \
        -v "$(shell pwd)":/go/src/${APP_PATH} \
        -w /go/src/${APP_PATH} \
        ${DOCKER_IMAGE} sh -c "make test"
else
	@go test -v $(shell glide novendor)
endif

test-coverage: depend
ifeq (${RUN},docker)
	@docker run --rm \
        -v "$(shell pwd)":/go/src/${APP_PATH} \
        -w /go/src/${APP_PATH} \
        ${DOCKER_IMAGE} sh -c 'make test-coverage'
else
	@echo "mode: set" > ./build/coverage.out; \
    for i in $$(go list ./... | grep -v "vendor"); do \
        go test -coverprofile=./build/cover.out $$i; \
        test -f ./build/cover.out && tail -n +2 ./build/cover.out >> ./build/coverage.out; \
    done; \
    go tool cover -html=./build/coverage.out -o ./build/coverage.html; \
    test -f ./build/cover.out && rm ./build/cover.out; \
    test -f ./build/coverage.out && rm ./build/coverage.out;
endif

test-report: depend
ifeq (${RUN},docker)
	@docker run --rm \
        -v "$(shell pwd)":/go/src/${APP_PATH} \
        -w /go/src/${APP_PATH} \
        ${DOCKER_IMAGE} sh -c "make test-report"
else
	@go test -v $(shell glide novendor) | go-junit-report > ./build/report.xml
endif