GOCMD=GO111MODULE=on CGO_ENABLED=0 go
GOBUILD=${GOCMD} build

# golangci-lint
LINTER := bin/golangci-lint

$(LINTER):
	curl -SL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v2.1.1

.PHONY: init
# Initialize environment
init:
	pre-commit install
	go install github.com/google/go-licenses@latest
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6
	go install github.com/google/wire/cmd/wire@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

.PHONY: generate
# Generate all
generate: proto api wire

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

.PHONY: run
# Run app
run:
	@if [ -z "$(app)" ]; then \
  	echo "error: app must exist. ex: app=account"; \
	else \
		echo "run: app/$(app)" && cd app/$(app) && $(MAKE) run; \
	fi

.PHONY: build
# Build app execute file. Use app=app_name to build specific service.
build:
	@if [ -z "$(app)" ]; then \
		find app -type d -depth 1 -print | xargs -L 1 bash -c 'echo "build: $$0" && cd "$$0" && $(MAKE) build'; \
	else \
		echo "build: app/$(app)" && cd app/$(app) && $(MAKE) build; \
	fi

.PHONY: start
# Start all app services. Use app=app_name to start specific service.
start: stop build
	@if [ -z "$(app)" ]; then \
		find app -type d -depth 1 -print | xargs -L 1 bash -c 'echo "start: $$0" && cd "$$0" && $(MAKE) start'; \
	else \
		echo "start: app/$(app)" && cd app/$(app) && $(MAKE) start; \
	fi

.PHONY: stop
# Stop all app services. Use app=app_name to stop specific service.
stop:
	@if [ -z "$(app)" ]; then \
		find app -type d -depth 1 -print | xargs -L 1 bash -c 'echo "stop: $$0" && cd "$$0" && $(MAKE) stop'; \
	else \
		echo "stop: app/$(app)" && cd app/$(app) && $(MAKE) stop; \
	fi

.PHONY: log
# Tail app service log file. Must use app=app_name to tail specific service.
log:
	@if [ -z "$(app)" ]; then \
  	echo "error: app must exist. ex: app=account"; \
	else \
		echo "log: app/$(app)" && cd app/$(app) && $(MAKE) log; \
	fi

.PHONY: test
# Run tests with race detector.
test:
	go test -race -cover ./...

.PHONY: vet
# Run go vet
vet:
	go vet ./...

.PHONY: license-check
# Run license check
license-check:
	go-licenses check ./...

.PHONY: lint
# Run lint
lint:
	golangci-lint run ./...

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
