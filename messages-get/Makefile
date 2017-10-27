all: fk-incoming.pb.go ingester

fk-incoming.pb.go: fk-incoming.proto
	protoc --go_out=./ fk-incoming.proto

ingester: *.go
	go build -o ingester *.go
