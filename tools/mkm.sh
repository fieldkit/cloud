#!/bin/bash

if ! command -v migrate &> /dev/null
then
    echo "golang-migrate could not be found. Install from https://github.com/golang-migrate/migrate"
    exit
fi

NAME=$1
if [ -z $NAME ]; then
    echo "Usage: mkm.sh <name>"
    exit 2
fi

migrate create -ext sql -dir migrations/primary $1
