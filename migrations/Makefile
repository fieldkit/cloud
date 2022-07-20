default: image

image:
	docker build -t conservify/fk-cloud-migrate .

dump:
	docker create --name temp conservify/fk-cloud-migrate
	docker export temp | tar t
	docker rm temp

migrate:
	docker run --rm --env 'MIGRATE_DATABASE_URL=${DATABASE_URL}' -v `pwd`:/migrations conservify/fk-cloud-migrate migrate
