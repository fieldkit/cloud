UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
GOARCH ?= amd64
GOOS ?= linux
endif

ifeq ($(UNAME_S),Darwin)
GOARCH ?= amd64
GOOS ?= darwin
endif

GO ?= env GOOS=$(GOOS) GOARCH=$(GOARCH) go
BUILD ?= $(abspath build)
WORKING_DIRECTORY ?= `pwd`
DOCKER_TAG ?= master

SERVER_SOURCES = $(shell find server -type f -name '*.go')

default: setup binaries tests gotests

ci: setup binaries tests

ci-db-tests:
	rm -rf active-schema
	mkdir -p active-schema
	docker stop fktests-pg || true
	docker rm fktests-pg || true
	docker run --name fktests-pg -e POSTGRES_DB=fieldkit -e POSTGRES_USER=fieldkit -e POSTGRES_PASSWORD=password -p "5432:5432" -d mdillon/postgis
	sleep 5 # TODO Add a loop here to check
	IP=`docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' fktests-pg`; \
	cd server && FIELDKIT_POSTGRES_URL="postgres://fieldkit:password@$$IP/fieldkit?sslmode=disable" go test -p 1 -v ./...
	docker stop fktests-pg || true

setup: legacy/src/js/secrets.js portal/src/secrets.js ocr-portal/src/js/secrets.js
	true || env GO111MODULE=on go get -u goa.design/goa/v3/...@v3

legacy/src/js/secrets.js: legacy/src/js/secrets.js.template
	cp $^ $@

ocr-portal/src/js/secrets.js: ocr-portal/src/js/secrets.js.template
	cp $^ $@

portal/src/secrets.js: portal/src/secrets.js.template
	cp $^ $@

binaries: $(BUILD)/server $(BUILD)/ingester $(BUILD)/fktool $(BUILD)/fkstreams

portal/node_modules:
	cd portal && npm install

tests: portal/node_modules
	cd portal && vue-cli-service test:unit

gotests:
	cd server && go test -p 1 ./...

dev-portal: portal/node_modules
	cd portal && npm run serve

all: binaries

server: $(BUILD)/server

ingester: $(BUILD)/ingester

fktool: $(BUILD)/fktool

$(BUILD)/server: $(SERVER_SOURCES)
	cd server/cmd/server && $(GO) build -o $@

$(BUILD)/ingester: $(SERVER_SOURCES)
	cd server/cmd/ingester && $(GO) build -o $@

$(BUILD)/fktool: server/tools/fktool/*.go $(SERVER_SOURCES)
	cd server/tools/fktool && $(GO) build -o $@ *.go

$(BUILD)/fkstreams: server/tools/fktool/*.go $(SERVER_SOURCES)
	cd server/tools/fkstreams && $(GO) build -o $@ *.go

generate:
	./tools/goa-generate.sh

generate-v3:
	cd server/api && goa gen github.com/fieldkit/cloud/server/api/design

generate-v3-example:
	mkdir -p server/api/example
	cd server/api/example && goa example github.com/fieldkit/cloud/server/api/design

clean:
	rm -rf $(BUILD)

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

.PHONY: schema-production

clean-postgres:
	docker-compose stop postgres
	rm -f active-schema

active-schema:
	ln -sf schema active-schema

restart-postgres: active-schema
	docker-compose stop postgres
	docker-compose rm -f postgres
	docker-compose up -d postgres

veryclean:

migrate-up:
	migrate -path migrations -database "postgres://fieldkit:password@127.0.0.1:5432/fieldkit?sslmode=disable&search_path=public" up

migrate-down:
	migrate -path migrations -database "postgres://fieldkit:password@127.0.0.1:5432/fieldkit?sslmode=disable&search_path=public" down

aws-image:
	cp ocr-portal/src/js/secrets.js.aws ocr-portal/src/js/secrets.js
	cp legacy/src/js/secrets.js.aws legacy/src/js/secrets.js
	cp portal/src/secrets.js.aws portal/src/secrets.js
	WORKING_DIRECTORY=$(WORKING_DIRECTORY) DOCKER_TAG=$(DOCKER_TAG) ./build.sh
