all: fk-data.proto.json fk-data.pb.go src/fk-data.pb.c src/fk-data.pb.h

node_modules/.bin/pbjs:
	npm install

fk-data.proto.json: node_modules/.bin/pbjs fk-data.proto
	pbjs fk-data.proto -t json -o fk-data.proto.json

src/fk-data.pb.c src/fk-data.pb.h: fk-data.proto
	protoc --nanopb_out=./src fk-data.proto

fk-data.pb.go: fk-data.proto
	protoc --go_out=./ fk-data.proto

veryclean:
