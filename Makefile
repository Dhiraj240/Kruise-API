GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

IMAGE_REPO = ryane

ifdef VERSION
	DOCKER_VERSION = $(VERSION)
	BINARY_VERSION = $(VERSION)
endif

ifndef VERSION
	ifneq ($(GIT_TAG),)
		DOCKER_VERSION = $(GIT_TAG)
		BINARY_VERSION = $(GIT_TAG)
	endif
endif

DOCKER_VERSION ?= ${GIT_SHA}
BINARY_VERSION ?= ${GIT_SHA}

LDFLAGS += -X deploy-wizard/cmd.version=${BINARY_VERSION}
LDFLAGS += -X deploy-wizard/cmd.gitCommit=${GIT_COMMIT}
LDFLAGS += -X deploy-wizard/cmd.gitTreeState=${GIT_DIRTY}

.PHONY: all
all: genrun

.PHONY: install
install:
	go install -ldflags '$(LDFLAGS)' ./cmd/server

.PHONY: build
build:
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o deploy-wizard -ldflags '$(LDFLAGS)' ./cmd/server

.PHONY: gen
gen:
	swagger generate server -t gen -f swagger.yaml --exclude-main -A deploy-wizard

.PHONY: run
run:
	go run cmd/server/main.go

genrun: gen run

.PHONY: test
test:
	@go test -parallel 4 ./...

.PHONY: docker
docker:
	docker build --force-rm --build-arg VERSION=${DOCKER_VERSION} -t ${IMAGE_REPO}/kruise-api:${DOCKER_VERSION} -f ./Dockerfile .

docker.push: docker
	docker push ${IMAGE_REPO}/kruise-api:${DOCKER_VERSION}
