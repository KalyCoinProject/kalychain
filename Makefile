
.PHONY: download-spec-tests
download-spec-tests:
	git submodule init
	git submodule update

.PHONY: bindata
bindata:
	go-bindata -pkg chain -o ./chain/chain_bindata.go ./chain/chains

.PHONY: protoc
protoc:
	protoc --go_out=. --go-grpc_out=. ./server/proto/*.proto
	protoc --go_out=. --go-grpc_out=. ./protocol/proto/*.proto
	protoc --go_out=. --go-grpc_out=. ./network/proto/*.proto
	protoc --go_out=. --go-grpc_out=. ./txpool/proto/*.proto
	protoc --go_out=. --go-grpc_out=. ./consensus/ibft/**/*.proto

.PHONY: build
build:
	$(eval LATEST_VERSION = $(shell git describe --tags --abbrev=0))
	$(eval COMMIT_HASH = $(shell git rev-parse --short HEAD))
	$(eval DATE = $(shell date +'%Y-%m-%d_%T'))
	go build -o kalychain -ldflags="-X 'github.com/KalyCoinProject/kalychain/versioning.Version=$(LATEST_VERSION)+$(COMMIT_HASH)+$(DATE)'" main.go

.PHONY: lint
lint:
	golangci-lint run -c lint-rule.yaml --timeout=2m

.PHONY: test
test: build
	PATH=$(shell pwd):${PATH} go test -count=1 -coverprofile coverage.out -timeout 28m ./...

.PHONY: generate-bsd-licenses
generate-bsd-licenses:
	./generate_dependency_licenses.sh BSD-3-Clause,BSD-2-Clause > ./licenses/bsd_licenses.json
