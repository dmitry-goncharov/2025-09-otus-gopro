BIN := "./bin/calendar"
DOCKER_IMG_APP="calendar-app:1.0.0"
DOCKER_IMG_SCHEDULER="calendar-scheduler:1.0.0"
DOCKER_IMG_SENDER="calendar-sender:1.0.0"
DOCKER_IMG_BASE="calendar-base:1.0.0"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

#Postrges
POSTGRES_USER ?= postgres
POSTGRES_PASSWORD ?= password
POSTGRES_DB ?= backend
POSTGRES_PORT ?= 5435
POSTGRES_CONTAINER := postgres-calendar

#Rabbit
RABBIT_USER ?= rabbit
RABBIT_PASSWORD ?= password
RABBIT_PORT ?= 5672
RABBIT_PORT_UI ?= 15672
RABBIT_CONTAINER := rabbitmq


run-postgres:
	docker run -d --name $(POSTGRES_CONTAINER) \
	-e POSTGRES_USER=$(POSTGRES_USER) \
	-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	-e POSTGRES_DB=$(POSTGRES_DB) \
	-p $(POSTGRES_PORT):5432 \
	-v postgres-data:/var/lib/postgresql/data \
	postgres:latest

run-rabbit:
	docker run -d --name $(RABBIT_CONTAINER) \
	-p $(RABBIT_PORT):5672 \
	-p $(RABBIT_PORT_UI):15672 \
	-e RABBITMQ_DEFAULT_USER=$(RABBIT_USER) \
	-e RABBITMQ_DEFAULT_PASS=$(RABBIT_PASSWORD) \
	rabbitmq:management

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

generate:
	go generate ./...

run: build
	$(BIN) -config ./configs/config.yaml

build-img-app:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_APP) \
		-f build/calendar/Dockerfile .

build-img-scheduler:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_SCHEDULER) \
		-f build/scheduler/Dockerfile .

build-img-sender:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_SENDER) \
		-f build/sender/Dockerfile .

build-img:
	build-img-app build-img-scheduler build-img-sender

build-img-base:
	docker build \
		-t $(DOCKER_IMG_BASE) \
		-f build/Dockerfile .

run-img: build-img-app
	docker run $(DOCKER_IMG_APP)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

test-integration: build-img-base build-img
	docker-compose -f ./deployments/integrationtest/docker-compose.yaml up --exit-code-from integrationtest && docker-compose -f ./deployments/integrationtest/docker-compose.yaml down -v

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.63.4

lint: install-lint-deps
	golangci-lint run ./...

up: build-img
	docker-compose -f ./deployments/docker-compose.yaml up -d

down:
	docker-compose -f ./deployments/docker-compose.yaml down -v

.PHONY: build run build-img build-img-base build-img-app build-img-scheduler build-img-sender run-img version test test-integration lint up down
