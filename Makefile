# This file is the development Makefile for the project in this repository.
# All variables listed here are used as substitution in these Makefile targets.

SERVICE-NAME = mcrawler

define ENV-CONFIGURATION
ENV='dev'
endef

################################################################################


# Install all dependencies required.
#
# NOTE: Docker & Docker Compose should already be installed.
.PHONY: install
install:
		GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.17.1

# Build project binaries.
.PHONY: build
build: lint
		cd cmd/ && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o mcrawler

# Build project Docker images for dev environment.
.PHONY: docker-build
docker-build: build
		cp cmd/mcrawler build/docker/mcrawler/ && \
		docker build -t timtosi/mcrawler:latest build/docker/mcrawler/

# Runs linter against the service codebase.
#
# NOTE: This rule require gcc to be found in the `$PATH`.
.PHONY: lint
lint:
		@golangci-lint run --config conf/golangci.yml && \
		echo "linter pass ok !"

# Runs test suite.
.PHONY: test
test: lint
		go test -v -coverprofile=coverage.txt -tags integration -race -cover -timeout=120s $$(glide novendor)

# Run project locally.
.PHONY: run
run:
		docker-compose -f deployments/docker-compose.yaml up
