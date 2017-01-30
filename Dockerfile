FROM golang:1.7.5-alpine
ENV CGO_ENABLED=0
ADD . /go/src/github.com/O-C-R/fieldkit
WORKDIR /go/src/github.com/O-C-R/fieldkit
RUN go install -ldflags '-extldflags "-static"' .
