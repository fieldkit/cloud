#!/bin/bash

set -xe

export FIELDKIT_ADDR=:8080
export FIELDKIT_POSTGRES_URL=postgres://fieldkit:password@127.0.0.1/fieldkit?sslmode=disable

make build/server

build/server
