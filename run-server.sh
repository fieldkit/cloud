#!/bin/bash

source aws.env

export FIELDKIT_ADDR=:8080
export FIELDKIT_POSTGRES_URL=postgres://fieldkit:password@127.0.0.1/fieldkit?sslmode=disable

export FIELDKIT_HTTP_SCHEME=http
export FIELDKIT_MEDIA_BUCKETS=fk-media,fkprod-media
export FIELDKIT_STREAMS_BUCKETS=fk-streams,fkprod-streams
export FIELDKIT_ARCHIVER=default

export FIELDKIT_INFLUX_DB_URL=http://127.0.0.1:8086
export FIELDKIT_INFLUX_DB_USERNAME=admin
export FIELDKIT_INFLUX_DB_PASSWORD=asdfasdf
export FIELDKIT_INFLUX_DB_ORG=fk
export FIELDKIT_INFLUX_DB_BUCKET=sensors
export FIELDKIT_INFLUX_DB_TOKEN=RxdlPKjlsec_NdsTabjga3qNxBU0nAAGEuUMSZVcWX-_R3VXi1YTAtku5fZ5cO8LcbooNIh6qrmjrZRZNqVAOQ==

export FIELDKIT_API_DOMAIN=api.fklocal.org:8080
export FIELDKIT_API_HOST=http://api.fklocal.org:8080
export FIELDKIT_API_DOMAIN=127.0.0.1:8080
export FIELDKIT_API_HOST=http://127.0.0.1:8080

set -xe

make setup server

mkdir -p .fs

build/server "$@"
