#!/bin/bash

source .env

export FIELDKIT_POSTGRES_URL="postgres://fieldkit:password@127.0.0.1/fieldkit?sslmode=disable"
export FIELDKIT_TIME_SCALE_URL="postgres://postgres:password@127.0.0.1:5433/fk?sslmode=disable"

export FIELDKIT_STREAMS_BUCKETS=fk-streams,fkprod-streams

# export FIELDKIT_INFLUX_DB_TOKEN="RxdlPKjlsec_NdsTabjga3qNxBU0nAAGEuUMSZVcWX-_R3VXi1YTAtku5fZ5cO8LcbooNIh6qrmjrZRZNqVAOQ=="
# export FIELDKIT_INFLUX_DB_URL="http://127.0.0.1:8086"
# export FIELDKIT_INFLUX_DB_ORG="fk"
# export FIELDKIT_INFLUX_DB_BUCKET="sensors"

make movedata && time build/movedata "$@"

