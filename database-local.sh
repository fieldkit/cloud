#!/bin/bash

if which psql; then
	psql -h 127.0.0.1 -U fieldkit -U fieldkit
else
	# docker run -it --rm postgres psql -h 127.0.0.1 -U fieldkit -U fieldkit
	docker exec -it cloud_postgres_1 psql -U fieldkit fieldkit
fi
