#!/bin/sh 

cd `dirname $0`

if [ ! -f active-schema ]; then
	ln -sf schema active-schema
fi

docker-compose down && docker-compose up --build
