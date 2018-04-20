all: test

install-tools:
	go get -u github.com/golang/dep/cmd/dep

install-deps:
	dep ensure

update-deps:
	dep ensure -update

test:
	go test
