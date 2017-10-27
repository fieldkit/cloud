#!/bin/sh

set -ex

cd `dirname $0`
docker-compose down && exec docker-compose up --build
