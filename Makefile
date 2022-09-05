# The import path is where your repository can be found.
# To import subpackages, always prepend the full import path.
# If you change this, run `make clean`. Read more: https://git.io/vM7zV
IMPORT_PATH := github.com/buro9/microcosm
GOCMD := go1.19

VERSION          := $(shell git describe --tags --always --dirty="-dev")
DATE             := $(shell date -u '+%Y-%m-%d-%H%M UTC')
VERSION_FLAGS    := -ldflags='-X "main.Version=$(VERSION)" -X "main.BuildTime=$(DATE)"'

# Space separated patterns of packages to skip in list, test, format.
IGNORED_PACKAGES := /vendor/

.PHONY: all
all: microcosm-web

.PHONY: microcosm-web
microcosm-web: .GOPATH/.ok
	$Q $(GOCMD) install -mod=mod -v $(VERSION_FLAGS) $(IMPORT_PATH)/cmd/microcosm-web

.PHONY: deps
deps:
	$Q $(GOCMD) list -m -u -mod=mod all
	$Q $(GOCMD) mod tidy
	$Q $(GOCMD) get -d -u ./...
	$Q $(GOCMD) mod vendor

run: microcosm-web
	$Q docker-compose up

refresh: microcosm-web
	$Q docker-compose stop web
	$Q docker-compose rm -f web
	$Q docker-compose up -d
