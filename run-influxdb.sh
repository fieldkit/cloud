#!/bin/bash

mkdir -p $PWD/influxdb

docker run -d -p 8086:8086 \
	--name influxdb \
	-v $PWD/influxdb/data:/var/lib/influxdb2 \
	-v $PWD/influxdb/config:/etc/influxdb2 \
	-e DOCKER_INFLUXDB_INIT_MODE=setup \
	-e DOCKER_INFLUXDB_INIT_USERNAME=admin \
	-e DOCKER_INFLUXDB_INIT_PASSWORD=asdfasdf \
	-e DOCKER_INFLUXDB_INIT_ORG=fk \
	-e DOCKER_INFLUXDB_INIT_BUCKET=sensors \
	-e DOCKER_INFLUXDB_INIT_ADMIN_TOKEN=RxdlPKjlsec_NdsTabjga3qNxBU0nAAGEuUMSZVcWX-_R3VXi1YTAtku5fZ5cO8LcbooNIh6qrmjrZRZNqVAOQ== \
	influxdb:2.2.0
