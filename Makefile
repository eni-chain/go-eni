BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%H')
APPNAME := eni

# don't override user values
ifeq (,$(VERSION))
  VERSION := $(shell git describe --exact-match 2>/dev/null)
  # if VERSION is empty, then populate it with branch's name and raw commit hash
  ifeq (,$(VERSION))
    VERSION := $(BRANCH)-$(COMMIT)
  endif
endif

# Update the ldflags with the app, client & server names
ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=$(APPNAME) \
	-X github.com/cosmos/cosmos-sdk/version.AppName=$(APPNAME)d \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'
build:
	@go build $(BUILD_FLAGS) -o build/$(APPNAME)d ./cmd/$(APPNAME)d
.PHONY: build

##############
###  Test  ###
##############

test-unit:
	@echo Running unit tests...
	@go test -mod=readonly -v -timeout 30m ./...

test-race:
	@echo Running unit tests with race condition reporting...
	@go test -mod=readonly -v -race -timeout 30m ./...

test-cover:
	@echo Running unit tests and creating coverage report...
	@go test -mod=readonly -v -timeout 30m -coverprofile=$(COVER_FILE) -covermode=atomic ./...
	@go tool cover -html=$(COVER_FILE) -o $(COVER_HTML_FILE)
	@rm $(COVER_FILE)

bench:
	@echo Running unit tests with benchmarking...
	@go test -mod=readonly -v -timeout 30m -bench=. ./...

test: govet govulncheck test-unit

.PHONY: test test-unit test-race test-cover bench

#################
###  Install  ###
#################

all: install

install:
	@echo "--> ensure dependencies have not been modified"
	@go mod verify
	@echo "--> installing $(APPNAME)d"
	@go install $(BUILD_FLAGS) -mod=readonly ./cmd/$(APPNAME)d

.PHONY: all install

##################
###  Protobuf  ###
##################

# Use this target if you do not want to use Ignite for generating proto files
GOLANG_PROTOBUF_VERSION=1.28.1
GRPC_GATEWAY_VERSION=1.16.0
GRPC_GATEWAY_PROTOC_GEN_OPENAPIV2_VERSION=2.20.0

proto-deps:
	@echo "Installing proto deps"
	@go install github.com/bufbuild/buf/cmd/buf@v1.50.0
	@go install github.com/cosmos/gogoproto/protoc-gen-gogo@latest
	@go install github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v$(GOLANG_PROTOBUF_VERSION)
	@go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v$(GRPC_GATEWAY_VERSION)
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v$(GRPC_GATEWAY_PROTOC_GEN_OPENAPIV2_VERSION)
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

proto-gen:
	@echo "Generating protobuf files..."
	@ignite generate proto-go --yes

.PHONY: proto-deps proto-gen

#################
###  Linting  ###
#################

golangci_lint_cmd=golangci-lint
golangci_version=v1.61.0

lint:
	@echo "--> Running linter"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version)
	@$(golangci_lint_cmd) run ./... --timeout 15m

lint-fix:
	@echo "--> Running linter and fixing issues"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version)
	@$(golangci_lint_cmd) run ./... --fix --timeout 15m

.PHONY: lint lint-fix

###################
### Development ###
###################

govet:
	@echo Running go vet...
	@go vet ./...

govulncheck:
	@echo Running govulncheck...
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@govulncheck ./...

.PHONY: govet govulncheck

reset-eni-node:
	@echo Resetting eni node...
	rm -rf eni-node && git checkout  eni-node

reset-multi-node:
	@echo Resetting multi eni node...
	rm -rf eni-nodes && git checkout eni-nodes

start4-node:
	@echo Starting 4 eni nodes...
	nohup ./build/enid start --home=./eni-nodes/node1 &> ./build/node1.log &
	nohup ./build/enid start --home=./eni-nodes/node2 --evm.http_port 1856 --evm.ws_port 1859 &> ./build/node2.log &
	nohup ./build/enid start --home=./eni-nodes/node3 --evm.http_port 1866 --evm.ws_port 1869 &> ./build/node3.log &
	nohup ./build/enid start --home=./eni-nodes/node4 --evm.http_port 1876 --evm.ws_port 1879 &> ./build/node4.log &
	@echo Done starting 4 eni nodes.

stop4-node:
	@echo Stopping 4 eni nodes...
	@pkill -f "enid start --home=./eni-nodes/node1"
	@pkill -f "enid start --home=./eni-nodes/node2"
	@pkill -f "enid start --home=./eni-nodes/node3"
	@pkill -f "enid start --home=./eni-nodes/node4"
	@echo Done stopping 4 eni nodes.

# build target
build-loadtest:
	@echo Building loadtest...
	go build -o build/loadtest loadtest/*.go
	@cp loadtest/config.json build/config.json
	@echo Done building loadtest.
# run target
run-loadtest: build-loadtest
	@echo Running loadtest...
	build/loadtest

# clean target
clean-loadtest:
	@echo Cleaning loadtest...
	rm -rf build/loadtest
	@echo Done cleaning loadtest.

#deploy erc20
deploy_erc20:
	@echo deploy erc20 contract to  to local nodes...
	./loadtest/contracts/deploy_erc20new.sh http://localhost:8545

#deploy erc721
deploy_erc721:
	@echo deploy erc721 contract to  to local nodes...
	./loadtest/contracts/deploy_erc721new.sh http://localhost:8545

.PHONY: all build run clean

start-prometheus-grafana-dashboard:
	@echo Starting prometheus and grafana dashboard...
	docker-compose -f docker/prometheus-grafana/docker-compose.yml up -d
	@echo Done starting prometheus and grafana dashboard.

stop-prometheus-grafana-dashboard:
	@echo Stopping prometheus and grafana dashboard...
	docker-compose -f docker/prometheus-grafana/docker-compose.yml down
	@echo Done stopping prometheus and grafana dashboard.