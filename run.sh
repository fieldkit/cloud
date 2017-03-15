#!/bin/sh
docker-compose down
docker-compose --verbose up --build
