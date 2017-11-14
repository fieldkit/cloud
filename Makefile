BUILD=build/

all: $(BUILD)/server $(BUILD)/db-tester $(BUILD)/sqs-worker $(BUILD)/sqs-sender

$(BUILD)/server: server/*.go
	go build -o $@ $^

$(BUILD)/db-tester: server/db-tester/*.go
	go build -o $@ $^

$(BUILD)/sqs-worker: server/sqs-worker/*.go
	go build -o $@ $^

$(BUILD)/sqs-sender: server/sqs-sender/*.go
	go build -o $@ $^

image:
	cd server && docker build -t conservify/fk-cloud-server .

generate:
	# https://github.com/goadesign/goa/issues/923
	mv server/vendor server/vendor-temp
	(cd $(GOPATH)/src/github.com/fieldkit/cloud/server && go generate) || true
	mv server/vendor-temp server/vendor

