GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION := $(shell git branch | grep \* | cut -d ' ' -f2)
BIN_NAME:=$(notdir $(shell pwd))
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
IMAGE_NAME := "stonejianbu/${BIN_NAME}"

.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	@echo "building ${BIN_NAME} ${VERSION}"
	go build ${ENABLE_RACE} -ldflags "-X main.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X main.Version=${VERSION}" -o bin/server main.go

.PHONY: package
package:
	@echo "building image ${BIN_NAME} ${VERSION} ${GIT_COMMIT}"
	docker build ${DOCKER_BUILD_ARGS} --push --build-arg APP_NAME=${BIN_NAME} --build-arg VERSION=${VERSION} --build-arg GIT_COMMIT=${GIT_COMMIT} -t ${IMAGE_NAME}:${VERSION} .