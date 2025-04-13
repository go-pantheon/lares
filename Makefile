GOCMD=GO111MODULE=on CGO_ENABLED=0 go
GOBUILD=${GOCMD} build

.PHONY: init
# Initialize environment
init:
	go install github.com/google/wire/cmd/wire@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest

.PHONY: generate
# Generate all
generate: proto api

.PHONY: version
# Show the generated version
version:
	@find app -type d -depth 1 -print | xargs -L 1 bash -c 'echo "version: $$0" && cd "$$0" && $(MAKE) version'

.PHONY: wire
# Generate wire
wire:
	@find app -type d -depth 1 -print | xargs -L 1 bash -c 'echo "wire: $$0" && cd "$$0" && $(MAKE) wire'

.PHONY: proto
# Generate internal proto struct
proto:
	buf generate

.PHONY: api
# Generate API
api:
	buf generate --template=buf.gen.server.yaml

.PHONY: build
# build execute file
build:
	@if [ -z "$(app)" ]; then \
		find app -type d -depth 1 -print | xargs -L 1 bash -c 'cd "$$0" && pwd && $(MAKE) build'; \
	else \
		cd app/$(app) && pwd && $(MAKE) build; \
	fi

.PHONY: run
# Start all project services
run: start log

.PHONY: log
# tail -f app/gate/bin/debug.log
log:
	@if [ -z "$(app)" ]; then \
  	    echo "error: app must exist. ex: app=player"; \
	else \
		cd app/$(app) && pwd && $(MAKE) log; \
	fi

.PHONY: start
# start all project services
start: stop build
	@if [ -z "$(app)" ]; then \
		find app -type d -depth 1 -print | xargs -L 1 bash -c 'cd "$$0" && pwd && $(MAKE) start'; \
	else \
		cd app/$(app) && pwd && $(MAKE) start; \
	fi

.PHONY: stop
# stop all project services
stop:
	@if [ -z "$(app)" ]; then \
		find app -type d -depth 1 -print | xargs -L 1 bash -c 'cd "$$0" && pwd && $(MAKE) stop'; \
	else \
		cd app/$(app) && pwd && $(MAKE) stop; \
	fi

.PHONY: docker
# build docker image
docker:
	@if [ -z "$(app)" ]; then \
		find app -type d -depth 1 -print | xargs -L 1 bash -c 'cd "$$0" && pwd && $(MAKE) docker'; \
	else \
		cd app/$(app) && pwd && $(MAKE) docker; \
	fi

.PHONY: test
# Run tests
test:
	go test -v ./... -cover

.PHONY: vet
# Run go vet
vet: 
	go vet ./...

# Show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
