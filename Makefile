out_binary=bin/orbiter
docker_image_fqdn=gianarb/orbiter
PACKAGES=$(shell go list ./... | grep -v /vendor/)
RACE=$(shell test $$(go env GOARCH) != "amd64" || (echo "-race"))
REV ?= $$(git rev-parse --short=7 HEAD)

.PHONY: build
build: clean $(out_binary)

$(out_binary): bin
	mkdir -p ./bin
	CGO_ENABLED=0 GOOS=linux go build -o $(out_binary) -a -tags netgo -ldflags '-w' .

bin:
	mkdir -p bin/

.PHONY: build/docker-image
docker-image:
	docker build -t $(docker_image_fqdn):latest .

docker-commit-publish: docker-image
	docker tag $(docker_image_fqdn):latest $(docker_image_fqdn):${REV}
	docker push $(docker_image_fqdn):${REV}

docker-latest-publish:
	docker push $(docker_image_fqdn):latest

.PHONY: clean
help: ## this help
	awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.PHONY: ci test vet
ci: vet test

test: ## run tests, except integration tests
	go test ${RACE} ${PACKAGES}

vet: ## run go vet
	test -z "$$(go vet ${PACKAGES} 2>&1 | grep -v '*composite literal uses unkeyed fields|exit status 0)' | tee /dev/stderr)"

.PHONY: clean
clean:
	rm -Rf bin/
