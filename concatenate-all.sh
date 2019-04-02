#!/bin/bash

ROOT=http://127.0.0.1:8080
ROOT=https://api.fkdev.org
DEVICES_JSON=/tmp/devices.json

curl -s $ROOT/files/devices > $DEVICES_JSON

for url in `cat $DEVICES_JSON | jq -r ".devices[].urls.data.generate" | head -n 20`; do
    echo $url
    curl -v $url
done


for url in `cat $DEVICES_JSON | jq -r ".devices[].urls.logs.generate" | head -n 20`; do
    echo $url
    curl -v $url
done
