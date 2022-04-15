#!/bin/bash

export FIELDKIT_INFLUX_DB_TOKEN="RxdlPKjlsec_NdsTabjga3qNxBU0nAAGEuUMSZVcWX-_R3VXi1YTAtku5fZ5cO8LcbooNIh6qrmjrZRZNqVAOQ=="

make intoflux && build/intoflux

