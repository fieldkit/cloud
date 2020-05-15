#!/bin/bash

docker exec -it cloud_postgres_1 pg_dump -U fieldkit fieldkit > fk.dump
