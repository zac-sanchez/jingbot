BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
GIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u +%Y-%m-%d-%H:%M:%S)

UNAME := $(shell uname)
ifeq ($(UNAME), Darwin)
  SED := gsed
else
  SED := sed
endif

# function to increment/add minor version, eg: v1.2.3 -> v1.2.4, or v1.2 -> v1.2.0
bumpMinor = $(shell echo $(1) | $(SED) -r 's/^(v?[0-9]+.[0-9]+.)(.+)/echo \1$$((\2+1))/ge;s/^(v?[0-9]+.[0-9]+)$$/\1.0/g')

CURRENT_VERSION := $(shell git describe --abbr=0 | $(SED) -r 's/^.*-//')  # nearest annotated tag
MANGLED_BRANCH :=  $(shell echo $(BRANCH) | tr -cd '[:alnum:]' | tr '[A-Z]' '[a-z]')
NEXT_VERSION := $(call bumpMinor,$(CURRENT_VERSION))

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test:
	go test -v ./...

.PHONY: local-docker-image
local-docker-image: cmd/jingbot/*.go ## Build local docker image"
	docker build --build-arg GIT_HASH=$(GIT_HASH) --build-arg BUILD_DATE=$(BUILD_DATE) --build-arg NEXT_VERSION=$(NEXT_VERSION) -t jingbot:$(GIT_HASH) -f Dockerfile .

.PHONY: docker-compose-run
docker-compose-run: ## Run jingbot in docker-compose accompanied by statsd
	@echo "You may need to run below command to build or rebuild docker images:\n make local-docker-image"
	GIT_HASH=$(GIT_HASH) IMAGE_NAME=jingbot docker-compose up

.PHONY: frontend
frontend:
	statik -src=$${PWD}/frontend/static -dest $${PWD}/pkg/web/

.PHONY: prod-docker-image
prod-docker-image: cmd/jingbot/*.go
	## Store your production build instructions in ./Makefile-local
	make -f Makefile-local prod-docker-image

.PHONY: prod-push
prod-push:
	## Store your production push instructions in ./Makefile-local
	make -f Makefile-local docker-push

.PHONY: get-next-version
get-next-version:
	@echo $(NEXT_VERSION)

.PHONY: clean
clean: ## Delete all the artifacts from makefile
	rm -rf artifacts
