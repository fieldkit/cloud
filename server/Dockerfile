FROM golang:1.8.1
ENV CGO_ENABLED=0
COPY . $GOPATH/src/github.com/fieldkit/cloud/server
WORKDIR $GOPATH/src/github.com/fieldkit/cloud/server
RUN go get ./... && go install -ldflags '-extldflags "-static"' .
