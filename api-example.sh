#!/bin/bash

AUTH=`curl -s -X POST -d '{ "username": "demo-user", "password": "asdfasdfasdf" }' -D - http://127.0.0.1:8080/login -o /dev/null | grep -Fi Authorization | tr -d '\r'`

echo "$AUTH"

curl -s -H "$AUTH" http://127.0.0.1:8080/projects/@/demo/expeditions/@/demo-expedition/inputs | python -m json.tool

curl -s -H "$AUTH" http://127.0.0.1:8080/projects/@/www/expeditions/@/www-expedition/inputs | python -m json.tool

curl -s -H "$AUTH" http://127.0.0.1:8080/projects/@/www/expeditions/@/www-expedition/inputs/devices | python -m json.tool

curl -s -H "$AUTH" http://127.0.0.1:8080/inputs/devices/2 | python -m json.tool

curl -s -H "$AUTH" -X POST -d '{ "name": "TEST", "key": "TEST-A" }' http://127.0.0.1:8080/expeditions/2/inputs/devices | python -m json.tool

curl -s -H "$AUTH" -X PATCH -d '{ "key": "ST", "json_schema": "{}", "active": true }' http://127.0.0.1:8080/inputs/devices/2/schemas | python -m json.tool

curl -s -H "$AUTH" -X PATCH -d '{ "key": "A2", "json_schema": "{}", "active": true }' http://127.0.0.1:8080/inputs/devices/2/schemas | python -m json.tool
