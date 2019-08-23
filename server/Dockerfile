FROM golang:1.12.9
ENV CGO_ENABLED=0

COPY . /app
WORKDIR /app

RUN mkdir /app/build && \
    go build -o /app/build/server -ldflags '-extldflags "-static"' server.go && \
    go build -o /app/build/ingester -ldflags '-extldflags "-static"' ingester.go
