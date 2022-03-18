#!/bin/bash

source aws.env

export FIELDKIT_ADDR=:8080
export FIELDKIT_POSTGRES_URL=postgres://fieldkit:password@127.0.0.1/fieldkit?sslmode=disable
export FIELDKIT_API_DOMAIN=127.0.0.1:8080
export FIELDKIT_API_HOST=http://127.0.0.1:8080
export FIELDKIT_HTTP_SCHEME=http
export FIELDKIT_MEDIA_BUCKETS=fk-media,fkprod-media
export FIELDKIT_STREAMS_BUCKETS=fk-streams,fkprod-streams
export FIELDKIT_ARCHIVER=default

set -xe

make setup server

mkdir -p .fs

build/server "$@"
