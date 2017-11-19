
BUILD=build

default: binaries

binaries: $(BUILD)/server $(BUILD)/db-tester $(BUILD)/sqs-worker $(BUILD)/sqs-sender $(BUILD)/fkcli

all: binaries bundles

bundles: frontend/build

frontend/build:
	cd frontend && yarn build

SERVER_SOURCES = $(shell find server -type f -name '*.go' -not -path "server/vendor/*")
FRONTEND_BUNDLE_SOURCES = $(shell find frontend/src/js -type f -name '*.js')

$(BUILD)/server: $(SERVER_SOURCES)
	go build -o $@ server/server.go

$(BUILD)/db-tester: server/db-tester/*.go $(SERVER_SOURCES)
	go build -o $@ server/db-tester/*.go

$(BUILD)/sqs-worker: server/sqs-worker/*.go $(SERVER_SOURCES)
	go build -o $@ server/sqs-worker/*.go

$(BUILD)/sqs-sender: server/sqs-sender/*.go $(SERVER_SOURCES)
	go build -o $@ server/sqs-sender/*.go

$(BUILD)/fkcli: server/api/tool/fieldkit-cli/*.go $(SERVER_SOURCES)
	go build -o $@ server/api/tool/fieldkit-cli/*.go

image:
	cd server && docker build -t conservify/fk-cloud-server .

generate:
	mv server/vendor server/vendor-temp # See https://github.com/goadesign/goa/issues/923
	(cd $(GOPATH)/src/github.com/fieldkit/cloud/server && go generate) || true
	mv server/vendor-temp server/vendor

clean:
	rm -rf build frontend/build
