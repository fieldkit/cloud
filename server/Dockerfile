FROM golang:1.19.1
ENV CGO_ENABLED=0

COPY . /app
WORKDIR /app

RUN mkdir /app/build

RUN cd cmd/movedata && go build -o /app/build/movedata -ldflags '-extldflags "-static"' *.go
RUN cd cmd/fktool && go build -o /app/build/fktool -ldflags '-extldflags "-static"' *.go
RUN cd cmd/webhook && go build -o /app/build/webhook -ldflags '-extldflags "-static"' *.go
RUN cd cmd/server && go build -o /app/build/server -ldflags '-extldflags "-static"' *.go
RUN cd cmd/ingester && go build -o /app/build/ingester -ldflags '-extldflags "-static"' *.go
