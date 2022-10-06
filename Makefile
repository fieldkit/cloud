VERSION_MAJOR = 0
VERSION_MINOR = 3
VERSION_PATCH = 29
VERSION_PREL ?= $(BUILD_NUMBER)
GIT_LOCAL_BRANCH ?= unknown
GIT_HASH ?= $(shell git log -1 --format=%h)
CI_CONTAINER_NAME ?= "fktests-$(GIT_LOCAL_BRANCH)"

VERSION := "$(VERSION_MAJOR).$(VERSION_MINOR).$(VERSION_PATCH)-$(GIT_LOCAL_BRANCH).$(VERSION_PREL)-$(GIT_HASH)"

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
JSPKG ?= yarn --mutex network
endif

GO ?= env GOOS=$(GOOS) GOARCH=$(GOARCH) go
BUILD ?= $(abspath build)
WORKING_DIRECTORY ?= $(shell pwd)
SERVER_SOURCES = $(shell find server -type f -name '*.go')
DOCKER_TAG ?= main

default: setup binaries jstests gotests charting-tests

setup: portal/src/secrets.ts

charting/node_modules:
	cd charting && $(JSPKG) install

charting-setup: charting/node_modules

charting-tests: charting-setup
	cd charting && tsc --esModuleInterop -t es5 server.ts

portal/src/secrets.ts: portal/src/secrets.ts.template
	cp $^ $@

binaries: $(BUILD)/server $(BUILD)/ingester $(BUILD)/fktool $(BUILD)/fkdata $(BUILD)/sanitizer $(BUILD)/webhook $(BUILD)/scratch $(BUILD)/movedata

portal/node_modules:
	cd portal && $(JSPKG) install

cycle-checks:
	npx madge --circular --extensions ts ./portal/src

jstests: portal/node_modules
	cd portal && $(JSPKG) install
	cd portal && node_modules/.bin/vue-cli-service build
	cd portal && node_modules/.bin/vue-cli-service test:unit

gotests:
	cd server && go test -p 1 -coverprofile=coverage.data ./...

gotest-specific:
	cd server && go test -p 1 github.com/fieldkit/cloud/server/api
	cd server && go test -p 1 github.com/fieldkit/cloud/server/data
	cd server && go test -p 1 github.com/fieldkit/cloud/server/ingester
	cd server && go test -p 1 github.com/fieldkit/cloud/server/backend

view-coverage:
	cd server && go tool cover -html=coverage.data

server: $(BUILD)/server

ingester: $(BUILD)/ingester

fktool: $(BUILD)/fktool

fkdata: $(BUILD)/fkdata

sanitizer: $(BUILD)/sanitizer

webhook: $(BUILD)/webhook

movedata: $(BUILD)/movedata

scratch: $(BUILD)/scratch

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

$(BUILD)/webhook: server/cmd/webhook/*.go $(SERVER_SOURCES)
	cd server/cmd/webhook && $(GO) build -o $@ *.go

$(BUILD)/movedata: server/cmd/movedata/*.go $(SERVER_SOURCES)
	cd server/cmd/movedata && $(GO) build -o $@ *.go

$(BUILD)/scratch: server/cmd/scratch/*.go $(SERVER_SOURCES)
	cd server/cmd/scratch && $(GO) build -o $@ *.go

generate:
	cd server/api && $(GOPATH)/bin/goa gen github.com/fieldkit/cloud/server/api/design

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
	@echo "CREATE DATABASE keycloak;" >> schema-production/0.sql
	@echo "UPDATE fieldkit.gue_jobs SET run_at = run_at + INTERVAL '10 year' WHERE run_at <= NOW();" >> schema-production/2.sql
	@for f in `find schema-production -name "*.xz"`; do                   \
		echo $$f                                                     ;\
		xz -d < $$f > schema-production/1.sql                        ;\
		rm $$f                                                       ;\
	done
	@for f in `find schema-production -name "*.bz2"`; do                  \
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
	sudo whoami
	docker-compose stop postgres
	docker-compose rm -f postgres
	sudo rm -rf .db
	docker-compose up -d postgres

veryclean:

run-server: server
	./run-server.sh

migrate-image:
	cd migrations && make image

migrate-up:
	cd migrations && MIGRATE_PATH=`pwd`/primary MIGRATE_DATABASE_URL="postgres://fieldkit:password@127.0.0.1:5432/fieldkit?sslmode=disable" go run main.go migrate

migrate-up-tsdb:
	cd migrations && MIGRATE_PATH=`pwd`/tsdb MIGRATE_DATABASE_URL="postgres://postgres:password@127.0.0.1:5433/fk?sslmode=disable" go run main.go migrate

ci: setup binaries jstests charting-setup

ci-db-tests:
	rm -rf active-schema
	mkdir -p active-schema
	@echo "DROP SCHEMA fieldkit;" > active-schema/0.sql
	docker stop $(CI_CONTAINER_NAME) || true
	docker rm $(CI_CONTAINER_NAME) || true
	docker network create docker_default || true
	docker run --name $(CI_CONTAINER_NAME) -e POSTGRES_DB=fieldkit -e POSTGRES_USER=fieldkit -e POSTGRES_PASSWORD=password --network=docker_default -d mdillon/postgis
	sleep 5 # TODO Add a loop here to check
	IP=`docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(CI_CONTAINER_NAME)`; \
	cd server && FIELDKIT_POSTGRES_URL="postgres://fieldkit:password@$$IP/fieldkit?sslmode=disable" go test -p 1 -v -coverprofile=coverage.data ./...
	docker stop $(CI_CONTAINER_NAME) || true

write-version:
	echo $(VERSION) > version.txt

docker-images:
	cp portal/src/secrets.ts.aws portal/src/secrets.ts
	WORKING_DIRECTORY=$(WORKING_DIRECTORY) DOCKER_TAG=$(DOCKER_TAG) VERSION=$(VERSION) ./build.sh

sanitize: sanitizer
	mkdir -p sanitize-data
	rsync -zvua schema-production/* sanitize-data
	docker run --rm --name fksanitize-pg -e POSTGRES_DB=fieldkit -e POSTGRES_USER=fieldkit -e POSTGRES_PASSWORD=password -p "5432:5432" -v `pwd`/sanitize-data:/docker-entrypoint-initdb.d/:ro -d postgis/postgis:12-2.5-alpine
	$(BUILD)/sanitizer --waiting
	docker exec fksanitize-pg pg_dump 'postgres://fieldkit:password@127.0.0.1/fieldkit?sslmode=disable' | bzip2 > db-sanitized.sql.bz2
	docker stop fksanitize-pg

.PHONY: schema-production sanitize
