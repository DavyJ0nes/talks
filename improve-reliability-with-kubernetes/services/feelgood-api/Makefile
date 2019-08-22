.DEFAULT_TARGET=help
.PHONY: all
all: help

# VARIABLES
USERNAME = awesome-corp
APP_NAME = feelgood-api

GO_VERSION ?= 1.13-rc
GO_PROJECT_PATH ?= github.com/davyj0nes/talks/productionise-using-sidecars/services/feelgood-api
GO_FILES = $(shell go list ./... | grep -v /vendor/)

APP_PORT = 80
LOCAL_PORT = 8080

VERSION = 0.0.1
COMMIT = $(shell git rev-parse HEAD | cut -c 1-6)
BUILD_TIME = $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')

BUILD_PREFIX = CGO_ENABLED=0 GOOS=linux
BUILD_FLAGS = -a -tags netgo --installsuffix netgo
LDFLAGS = -ldflags "-s -w -X ${GO_PROJECT_PATH}/cmd.Release=${VERSION} -X ${GO_PROJECT_PATH}/cmd.Commit=${COMMIT} -X ${GO_PROJECT_PATH}/cmd.BuildTime=${BUILD_TIME}"
DOCKER_GO_BUILD = docker run --rm -v "$(GOPATH)":/go -v "$(CURDIR)":/go/src/app -w /go/src/app golang:${GO_VERSION}
GO_BUILD_STATIC = $(BUILD_PREFIX) go build $(BUILD_FLAGS) $(LDFLAGS)

DOCKER_RUN_CMD = docker run -it --rm --name ${APP_NAME} ${USERNAME}/${APP_NAME}:${VERSION} "\$$@"

# COMMANDS

## binary: builds a statically linked binary of the application (used in Docker image)
.PHONY: binary
binary:
	$(call blue, "# Building Golang Binary...")
	@${DOCKER_GO_BUILD} sh -c 'go get && ${GO_BUILD_STATIC} -o ${APP_NAME} main.go message.go'

## image: builds a docker image for the application
.PHONY: image
image: binary
	$(call blue, "# Building Docker Image...")
	@docker build --no-cache -t ${USERNAME}/${APP_NAME}:${VERSION} .
	@docker tag ${USERNAME}/${APP_NAME}:${VERSION} ${USERNAME}/${APP_NAME}:latest
	@$(MAKE) clean

## publish: pushes the tagged docker image to docker hub
.PHONY: publish
publish: image
	$(call blue, "# Publishing Docker Image...")
	@docker push docker.io/${USERNAME}/${APP_NAME}:${VERSION}

## run: runs the application locally
.PHONY: run
run:
	$(call blue, "# Running App...")
	@docker run -it --rm -v "$(GOPATH)":/go -v "$(CURDIR)":/go/src/app -p ${LOCAL_PORT}:${APP_PORT} -w /go/src/app golang:${GO_VERSION} go run cmd/main.go

## run_image: builds and runs the docker image locally
.PHONY: run_image
run_image: image
	$(call blue, "# Running Docker Image Locally...")
	@docker run -it --rm --name ${APP_NAME} -p ${LOCAL_PORT}:${APP_PORT} ${USERNAME}/${APP_NAME}:${VERSION}

## test: run test suitde for application
.PHONY: test
test:
	$(call blue, "# Testing Golang Code...")
	@docker run --rm -it -v "$(GOPATH):/go" -v "$(CURDIR)":/go/src/app -w /go/src/app golang:${GO_VERSION} sh -c 'go test -v -race ${GO_FILES}' 

## clean: remove binary from non release directory
.PHONY: clean
clean: 
	@rm -f ${APP_NAME} 

## help: Show this help message
.PHONY: help
help: Makefile
	@echo "${APP_NAME} - v${VERSION}"
	@echo
	@echo " Choose a command run in "$(APP_NAME)":"
	@echo
	@sed -n 's/^## //p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

# FUNCTIONS
define blue
	@tput setaf 4
	@echo $1
	@tput sgr0
endef
