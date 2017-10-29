#!/bin/bash

docker exec -it ingestion_postgres_1 psql -U fieldkit fieldkit
