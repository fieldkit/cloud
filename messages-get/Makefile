all: fk-incoming.pb.go

fk-incoming.pb.go: fk-incoming.proto
	protoc --go_out=./ fk-incoming.proto
