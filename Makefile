BUILD=$(shell pwd)/build
dirs=src sqs-tester sqs-worker

all:
	mkdir -p $(BUILD)
	set -xe; for dir in $(dirs); do       \
		(cd $$dir && BUILD=$(BUILD) make);  \
	done

clean:
	rm -rf $(BUILD)


run: all env
	$(BUILD)/ingester

server:
	docker-compose down && docker-compose up --build

deps:
	go get github.com/golang/protobuf/proto
	go get github.com/lib/pq
	go get github.com/aws/aws-sdk-go/aws
