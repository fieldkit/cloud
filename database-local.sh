#!/bin/bash

docker exec -i ingestion_postgres_1 psql -U fieldkit fieldkit
