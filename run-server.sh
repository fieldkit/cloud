#!/bin/bash

source .env

if [ -z "$FIELDKIT_DB_HOST" ]; then
	export FIELDKIT_DB_HOST=127.0.0.1
fi

export FIELDKIT_ADDR=:8080
export FIELDKIT_POSTGRES_URL=postgres://fieldkit:password@${FIELDKIT_DB_HOST}/fieldkit?sslmode=disable
export FIELDKIT_TIME_SCALE_URL=postgres://postgres:password@${FIELDKIT_DB_HOST}:5433/fk?sslmode=disable

export FIELDKIT_HTTP_SCHEME=http
export FIELDKIT_MEDIA_BUCKETS=fk-media,fkprod-media
export FIELDKIT_STREAMS_BUCKETS=fk-streams,fkprod-streams
export FIELDKIT_ARCHIVER=default

export FIELDKIT_WORKERS=5

# Most of the time, development should be fine using this when jumping from
# working on a local machine to accessing that machine over the LAN for testing
# mobile.
export FIELDKIT_API_DOMAIN="{_:127\.0\.0\.1|192\.168\.0\.100}:8080"


# This may need adjustment for testing scenarios where the server returns a URL.
export FIELDKIT_API_HOST=http://127.0.0.1:8080

set -xe

make setup server

mkdir -p .fs

build/server "$@"
