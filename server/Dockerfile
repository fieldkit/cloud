FROM golang:1.14
ENV CGO_ENABLED=0

COPY . /app
WORKDIR /app

RUN mkdir /app/build && \
    go build -o /app/build/server -ldflags '-extldflags "-static"' cmd/server/server.go && \
    go build -o /app/build/ingester -ldflags '-extldflags "-static"' cmd/ingester/ingester.go
