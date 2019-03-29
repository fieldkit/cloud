#!/bin/bash

ROOT=http://127.0.0.1:8080
ROOT=https://api.fkdev.org

for url in `curl -s $ROOT/files/devices | jq -r ".devices[].urls.data.generate"`; do
    echo $url
    # curl -v $url
done
