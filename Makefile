GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION := v1.0.0
BIN_NAME:=$(notdir $(shell pwd))
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
IMAGE_NAME := "stonejianbu/${BIN_NAME}"

.PHONY: run
run:
	go run main.go

up:
	docker compose up -d

down:
	docker compose down

restart:
	docker compose restart

.PHONY: build
build:
	@echo "building ${BIN_NAME} ${VERSION}"
	go build ${ENABLE_RACE} -ldflags "-X main.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X main.Version=${VERSION}" -o bin/server main.go

.PHONY: package
package:
	@echo "building image ${BIN_NAME} ${VERSION} ${GIT_COMMIT}"
	docker build  -t ${IMAGE_NAME}:${VERSION} .