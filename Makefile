PACKAGES=$(shell go list ./... | grep -v /vendor/)
RACE := $(shell test $$(go env GOARCH) != "amd64" || (echo "-race"))

test: ## run tests, except integration tests
	@go test ${RACE} ${PACKAGES}

binaries:
	@echo "Compiling..."
	@mkdir -p ./bin
	@go build -i -o ./bin/orbiter
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -i -o ./bin/orbiter
	@echo "All done! The binaries is in ./bin let's have fun!"

build/docker: binaries
	@docker build -t gianarb/orbiter:latest .

vet: ## run go vet
	@test -z "$$(go vet ${PACKAGES} 2>&1 | grep -v '*composite literal uses unkeyed fields|exit status 0)' | tee /dev/stderr)"

help: ## this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

ci: vet test
