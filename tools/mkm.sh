#!/bin/bash

# TODO Check for migrate tool.

NAME=$1
if [ -z $NAME ]; then
    echo "Usage: mkm.sh <name>"
    exit 2
fi

migrate create -ext sql -dir migrations $1
