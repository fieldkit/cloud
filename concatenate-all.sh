#!/bin/bash

ROOT=http://127.0.0.1:8080
ROOT=https://api.fkdev.org

DEVICES_JSON=$(curl -s $ROOT/files/devices)

function pause_for_empty_queue() {
    while true; do
        queue_length=$(curl -s $ROOT/files/status | jq '.queued')

        echo "Queue:" $queue_length

        if [ $queue_length -lt 1 ]; then
            return
        fi

        echo "Waiting for empty queue..."

        sleep 10
    done
}

for device in `echo $DEVICES_JSON | jq -c '.devices | sort_by(.number_of_files) | .[] | { files: .number_of_files, data: .urls.data.generate, logs: .urls.logs.generate }'`; do
    data=$(echo $device | jq -r '.data')
    logs=$(echo $device | jq -r '.logs')

    pause_for_empty_queue

    echo $device
    curl $data
    curl $logs
done

