BUILD=build

all: $(BUILD)/server $(BUILD)/db-tester $(BUILD)/sqs-worker $(BUILD)/sqs-sender $(BUILD)/fkcli

SERVER_SOURCES = $(shell find server -type f -name '*.go' -not -path "server/vendor/*")

$(BUILD)/server: $(SERVER_SOURCES)
	go build -o $@ server/server.go

$(BUILD)/db-tester: server/db-tester/*.go $(SERVER_SOURCES)
	go build -o $@ $^

$(BUILD)/sqs-worker: server/sqs-worker/*.go $(SERVER_SOURCES)
	go build -o $@ $^

$(BUILD)/sqs-sender: server/sqs-sender/*.go $(SERVER_SOURCES)
	go build -o $@ $^

$(BUILD)/fkcli: server/api/tool/fieldkit-cli/*.go $(SERVER_SOURCES)
	go build -o $@ $^

image:
	cd server && docker build -t conservify/fk-cloud-server .

generate:
	mv server/vendor server/vendor-temp # See https://github.com/goadesign/goa/issues/923
	(cd $(GOPATH)/src/github.com/fieldkit/cloud/server && go generate) || true
	mv server/vendor-temp server/vendor

clean:
	rm -rf build
