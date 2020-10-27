UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
GOARCH ?= amd64
GOOS ?= linux
endif

ifeq ($(UNAME_S),Darwin)
GOARCH ?= amd64
GOOS ?= darwin
endif

ifeq (, $(shell which yarn))
JSPKG ?= npm
else
JSPKG ?= yarn
endif

GO ?= env GOOS=$(GOOS) GOARCH=$(GOARCH) go
BUILD ?= $(abspath build)
WORKING_DIRECTORY ?= $(shell pwd)
SERVER_SOURCES = $(shell find server -type f -name '*.go')
DOCKER_TAG ?= main

default: setup binaries jstests gotests

setup: portal/src/secrets.ts

portal/src/secrets.ts: portal/src/secrets.ts.template
	cp $^ $@

binaries: $(BUILD)/server $(BUILD)/ingester $(BUILD)/fktool $(BUILD)/fkdata $(BUILD)/sanitizer

portal/node_modules:
	cd portal && $(JSPKG) install

cycle-checks:
	npx madge --circular --extensions ts ./portal/src

jstests: portal/node_modules
	cd portal && $(JSPKG) install
	cd portal && npm run build
	cd portal && npm run test:unit

gotests:
	cd server && go test -p 1 -coverprofile=coverage.data ./...

view-coverage:
	cd server && go tool cover -html=coverage.data

server: $(BUILD)/server

ingester: $(BUILD)/ingester

fktool: $(BUILD)/fktool

fkdata: $(BUILD)/fkdata

sanitizer: $(BUILD)/sanitizer

$(BUILD)/server: $(SERVER_SOURCES)
	cd server/cmd/server && $(GO) build -o $@

$(BUILD)/ingester: $(SERVER_SOURCES)
	cd server/cmd/ingester && $(GO) build -o $@

$(BUILD)/fktool: server/cmd/fktool/*.go $(SERVER_SOURCES)
	cd server/cmd/fktool && $(GO) build -o $@ *.go

$(BUILD)/fkdata: server/cmd/fkdata/*.go $(SERVER_SOURCES)
	cd server/cmd/fkdata && $(GO) build -o $@ *.go

$(BUILD)/sanitizer: server/cmd/sanitizer/*.go $(SERVER_SOURCES)
	cd server/cmd/sanitizer && $(GO) build -o $@ *.go

generate:
	cd server/api && goa gen github.com/fieldkit/cloud/server/api/design

clean:
	rm -rf $(BUILD)
	rm -rf $(WORKING_DIRECTORY)/portal/node_modules

schema-production:
	@mkdir -p schema-production
	@rm -f schema-production/*.sql
	@rm -f active-schema
	@ln -sf schema-production active-schema
	@echo "CREATE USER fk;" > schema-production/0.sql
	@echo "CREATE USER server;" >> schema-production/0.sql
	@echo "CREATE USER rdsadmin;" >> schema-production/0.sql
	@for f in `find schema-production -name "*.bz2"`; do              \
		echo $$f                                                     ;\
		bunzip2 < $$f > schema-production/1.sql                      ;\
		rm $$f                                                       ;\
	done
	$(MAKE) restart-postgres

clean-postgres:
	docker-compose stop postgres
	rm -f active-schema

active-schema:
	mkdir -p schema
	ln -sf schema active-schema

restart-postgres: active-schema
	docker-compose stop postgres
	docker-compose rm -f postgres
	docker-compose up -d postgres

veryclean:

run-server: server
	./run-server.sh

migrate-up:
	migrate -path migrations -database "postgres://fieldkit:password@127.0.0.1:5432/fieldkit?sslmode=disable&search_path=public" up

migrate-down:
	migrate -path migrations -database "postgres://fieldkit:password@127.0.0.1:5432/fieldkit?sslmode=disable&search_path=public" down

ci: setup binaries jstests

ci-db-tests:
	rm -rf active-schema
	mkdir -p active-schema
	docker stop fktests-pg || true
	docker rm fktests-pg || true
	docker network create docker_default || true
	docker run --name fktests-pg -e POSTGRES_DB=fieldkit -e POSTGRES_USER=fieldkit -e POSTGRES_PASSWORD=password --network=docker_default -d mdillon/postgis
	sleep 5 # TODO Add a loop here to check
	IP=`docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' fktests-pg`; \
	cd server && FIELDKIT_POSTGRES_URL="postgres://fieldkit:password@$$IP/fieldkit?sslmode=disable" go test -p 1 -v -coverprofile=coverage.data ./...
	docker stop fktests-pg || true

aws-image:
	cp portal/src/secrets.ts.aws portal/src/secrets.ts
	WORKING_DIRECTORY=$(WORKING_DIRECTORY) DOCKER_TAG=$(DOCKER_TAG) ./build.sh

sanitize: sanitizer
	mkdir -p sanitize-data
	rsync -zvua schema-production/* sanitize-data
	docker run --rm --name fksanitize-pg -e POSTGRES_DB=fieldkit -e POSTGRES_USER=fieldkit -e POSTGRES_PASSWORD=password -p "5432:5432" -v `pwd`/sanitize-data:/docker-entrypoint-initdb.d/:ro -d mdillon/postgis
	$(BUILD)/sanitizer
	docker exec fksanitize-pg pg_dump 'postgres://fieldkit:password@127.0.0.1/fieldkit?sslmode=disable' | bzip2 > db-sanitized.sql.bz2
	docker stop fksanitize-pg

.PHONY: schema-production sanitize
