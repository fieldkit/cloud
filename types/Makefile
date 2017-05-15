SRC_DIR=./proto

proto: protogo protojs

protogo: go/document/fieldkit
	protoc -I=$(SRC_DIR)/document/fieldkit --go_out=import_prefix=github.com/O-C-R/fieldkit/types/go/document:go/document/fieldkit $(SRC_DIR)/document/fieldkit/*.proto

protojs: js
	protoc -I=$(SRC_DIR)/document/fieldkit --flow_out=js $(SRC_DIR)/document/fieldkit/*.proto

go/document/fieldkit:
	mkdir -p go/document/fieldkit

js:
	mkdir -p js

clean:
	rm -rf go/document js

.PHONY: proto
