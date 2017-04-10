FROM golang:1.8.1
ENV CGO_ENABLED=0
COPY . $GOPATH/src/github.com/O-C-R/fieldkit/server
WORKDIR $GOPATH/src/github.com/O-C-R/fieldkit/server
RUN go install -ldflags '-extldflags "-static"' .
