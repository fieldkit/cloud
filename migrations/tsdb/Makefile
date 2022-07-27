default: help

help:
	echo "See migrate target"

migrate:
	docker run --rm --env 'MIGRATE_DATABASE_URL=${DATABASE_URL}' -v `pwd`:/migrations conservify/fk-cloud-migrate migrate
