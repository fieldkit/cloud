all: build/send build/worker

build:
	mkdir -p build

build/send: build send.go
	go build -o build/send send.go

build/worker: build sqs-worker.go main.go
	go build -o build/worker main.go sqs-worker.go

clean:
	rm -f build
