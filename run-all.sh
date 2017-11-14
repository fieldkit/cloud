#!/bin/sh 

cd `dirname $0`

docker-compose down && docker-compose up --build
