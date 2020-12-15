#!/bin/bash

if [ -d ~/conservify/dev-ops/amis/services/auth-keycloak/theme ]; then
	docker run --rm --name keycloak \
		--link cloud_postgres_1 --network cloud_default \
		-p 9990:9990 \
		-p 8090:8080 \
		-e KEYCLOAK_USER=admin -e KEYCLOAK_PASSWORD=admin \
		-e DB_VENDOR=postgres -e DB_USER=fieldkit -e DB_PASSWORD=password -e DB_ADDR=cloud_postgres_1:5432 \
		-v ~/conservify/dev-ops/amis/services/auth-keycloak/theme:/opt/jboss/keycloak/themes/fk \
		conservify/fk-cloud-keycloak:active
fi
