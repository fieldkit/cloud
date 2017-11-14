#!/bin/bash

curl -s -X POST -d '{ "username": "demo-user", "password": "asdfasdfasdf" }' -D - http://127.0.0.1:8080/login -o /dev/null

AUTH="Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTA2Nzc1MjUsImlhdCI6MTUxMDY3MzkyNSwicmVmcmVzaF90b2tlbiI6ImEyOWw0ZW8yTUlFSXlnLTlrWEpiakNOTjY3ZyIsInNjb3BlcyI6ImFwaTphY2Nlc3MiLCJzdWIiOjEsInVzZXJuYW1lIjoiZGVtby11c2VyIn0.yCCr3hXaARyuHLsFhvVr-9Uvo-oUuPwaW28RaUtdEP_T10x9yxddfaaGUJno8c1WgdG0CoqtB5E9S3H1aMB-eQ"

curl -s -H "Authorization: $AUTH" http://127.0.0.1:8080/projects/@/demo/expeditions/@/demo-expedition/inputs
curl -s -H "Authorization: $AUTH" http://127.0.0.1:8080/projects/@/www/expeditions/@/www-expedition/inputs
