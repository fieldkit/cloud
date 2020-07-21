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
	docker run --name fktests-pg -e POSTGRES_DB=fieldkit -e POSTGRES_USER=fieldkit -e POSTGRES_PASSWORD=password --network=docker_default -d mdillon/postgis
	sleep 5 # TODO Add a loop here to check
	IP=`docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' fktests-pg`; \
	cd server && FIELDKIT_POSTGRES_URL="postgres://fieldkit:password@$$IP/fieldkit?sslmode=disable" go test -p 1 -v -coverprofile=coverage.data ./...
	docker stop fktests-pg || true

setup: portal/src/secrets.ts
	true || env GO111MODULE=on go get -u goa.design/goa/v3/...@v3

portal/src/secrets.ts: portal/src/secrets.ts.template
	cp $^ $@

binaries: $(BUILD)/server $(BUILD)/ingester $(BUILD)/fktool $(BUILD)/fkdata

portal/node_modules:
	cd portal && npm install

tests: portal/node_modules
	if [ -f portal/node_modules/.yarn-integrity ]; then  \
	cd portal && yarn install;                           \
	else                                                 \
	cd portal && npm install;                            \
	fi
	cd portal && vue-cli-service test:unit

gotests:
	cd server && go test -p 1 -coverprofile=coverage.data ./...

view-coverage:
	cd server && go tool cover -html=coverage.data

dev-portal: portal/node_modules
	cd portal && npm run serve

all: binaries

server: $(BUILD)/server

ingester: $(BUILD)/ingester

fktool: $(BUILD)/fktool

fkdata: $(BUILD)/fkdata

$(BUILD)/server: $(SERVER_SOURCES)
	cd server/cmd/server && $(GO) build -o $@

$(BUILD)/ingester: $(SERVER_SOURCES)
	cd server/cmd/ingester && $(GO) build -o $@

$(BUILD)/fktool: server/cmd/fktool/*.go $(SERVER_SOURCES)
	cd server/cmd/fktool && $(GO) build -o $@ *.go

$(BUILD)/fkdata: server/cmd/fkdata/*.go $(SERVER_SOURCES)
	cd server/cmd/fkdata && $(GO) build -o $@ *.go

generate:
	cd server/api && goa gen github.com/fieldkit/cloud/server/api/design

generate-goa-example:
	mkdir -p server/api/example
	cd server/api/example && goa example github.com/fieldkit/cloud/server/api/design

clean:
	rm -rf $(BUILD)
	rm -rf portal/node_modules

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
	cp portal/src/secrets.ts.aws portal/src/secrets.ts
	WORKING_DIRECTORY=$(WORKING_DIRECTORY) DOCKER_TAG=$(DOCKER_TAG) ./build.sh
