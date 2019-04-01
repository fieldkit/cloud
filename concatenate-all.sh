#!/bin/bash

ROOT=https://api.fkdev.org
ROOT=http://127.0.0.1:8080

for url in `curl -s $ROOT/files/devices | jq -r ".devices[].urls.data.generate" | head -n 10`; do
    echo $url
    curl -v $url
done
