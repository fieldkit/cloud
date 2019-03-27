#!/bin/bash


for url in `curl -s http://127.0.0.1:8080/files/devices | jq -r ".devices[].urls.data"`; do
    echo $url
    curl -v $url
done
