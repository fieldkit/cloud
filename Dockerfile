FROM node:18.17.0 AS node
ENV PUBLIC_PATH /
WORKDIR /app

COPY ./portal/package*.json ./
RUN npm install

COPY ./portal/ ./
RUN npm run build

RUN for e in css csv html js json map svg txt; do find . -iname '*.$e' -exec gzip -k9 {} \; ; done

FROM golang:latest AS golang
WORKDIR /app

COPY ./server/ ./

RUN mkdir build
RUN cd cmd/server && go build -o /app/build/server -ldflags '-extldflags "-static"' *.go
RUN cd cmd/ingester && go build -o /app/build/ingester -ldflags '-extldflags "-static"' *.go

FROM alpine:latest AS env
ARG GIT_HASH=missing
ARG VERSION=missing
WORKDIR /app

RUN echo "export GIT_HASH=$GIT_HASH" > static.env
RUN echo "export FIELDKIT_VERSION=$VERSION" >> static.env
RUN cat static.env

FROM scratch
COPY --from=golang /app/build/server /
COPY --from=golang /app/build/ingester /
COPY --from=node /app/build /portal
COPY --from=golang /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=env /app/static.env /etc/

# Downstream Dockerfile's require this. I'd like to use the above paths, though.
COPY --from=golang /etc/ssl/certs/ca-certificates.crt /
COPY --from=env /app/static.env /

EXPOSE 80
ENTRYPOINT ["/server"]