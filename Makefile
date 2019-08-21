modules = server/sqs-worker server/sqs-sender server/tools/fktool

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

SERVER_SOURCES = $(shell find server -type f -name '*.go' -not -path "server/vendor/*")
TESTING_SOURCES = $(shell find testing -type f -name '*.go' -not -path "server/vendor/*")

default: setup binaries

setup: fieldkit.env legacy/src/js/secrets.js portal/src/js/secrets.js server/inaturalist/secrets.go

fieldkit.env:
	echo FIELDKIT_ADDR=0.0.0.0:8080 > $@

legacy/src/js/secrets.js: legacy/src/js/secrets.js.template
	cp $^ $@

portal/src/js/secrets.js: portal/src/js/secrets.js.template
	cp $^ $@

server/inaturalist/secrets.go: server/inaturalist/secrets.go.template
	cp $^ $@

binaries: $(BUILD)/server $(BUILD)/sqs-worker $(BUILD)/sqs-sender $(BUILD)/fktool $(BUILD)/fkflash $(BUILD)/testing-random $(BUILD)/weather-proxy $(BUILD)/inaturalist

all: binaries

tests:
	cd legacy && yarn run flow

server: $(BUILD)/server

fktool: $(BUILD)/fktool

$(BUILD)/server: $(SERVER_SOURCES) server/inaturalist/secrets.go
	cd server && $(GO) build -o $@ server.go

$(BUILD)/fktool: server/tools/fktool/*.go $(SERVER_SOURCES) $(TESTING_SOURCES)
	cd server/tools/fktool && $(GO) build -o $@ *.go

$(BUILD)/fkflash: server/tools/fkflash/*.go $(SERVER_SOURCES) $(TESTING_SOURCES)
	cd server/tools/fkflash && $(GO) build -o $@ *.go

$(BUILD)/inaturalist: server/tools/inaturalist/*.go $(SERVER_SOURCES) $(TESTING_SOURCES) server/inaturalist/secrets.go
	cd server/tools/inaturalist && $(GO) build -o $@ *.go

$(BUILD)/sqs-worker: server/sqs-worker/*.go $(SERVER_SOURCES)
	cd server/sqs-worker && $(GO) build -o $@ *.go

$(BUILD)/sqs-sender: server/sqs-sender/*.go $(SERVER_SOURCES)
	cd server/sqs-sender && $(GO) build -o $@ *.go

$(BUILD)/testing-random: testing/random/*.go $(SERVER_SOURCES) $(TESTING_SOURCES)
	cd testing/random && $(GO) build -o $@ *.go

$(BUILD)/weather-proxy: testing/weather-proxy/*.go $(SERVER_SOURCES) $(TESTING_SOURCES)
	cd testing/weather-proxy && $(GO) build -o $@ *.go

install: all
	cp $(BUILD)/fktool $(INSTALLDIR)
	cp $(BUILD)/testing-random $(INSTALLDIR)
	cp $(BUILD)/sqs-sender $(INSTALLDIR)
	cp $(BUILD)/sqs-worker $(INSTALLDIR)
	@for d in $(modules); do                           \
		(cd $$d && echo $$d && go install) || exit 1;  \
	done

generate:
	./goa-generate.sh

deps: server/inaturalist/secrets.go
	cd server && go get ./...

clean:
	rm -rf $(BUILD)

clean-production:
	rm -rf schema-production

refresh-production: clean-production clone-production

schema-production:
	mkdir schema-production
	@if [ -d ~/conservify/dev-ops ]; then                                     \
		(cd ~/conservify/dev-ops/provisioning && ./db-dump.sh);                 \
		cp ~/conservify/dev-ops/schema.sql schema-production/000001.sql;        \
		cp ~/conservify/dev-ops/data.sql schema-production/000100.sql;          \
		cp schema/00000?.sql schema-production/;                                \
	else                                                                      \
		echo "No dev-ops directory found";                                      \
	fi

clone-production: schema-production
	rm -f active-schema
	ln -sf schema-production active-schema

clean-postgres:
	docker-compose stop postgres
	rm -f active-schema

active-schema:
	ln -sf schema active-schema

restart-postgres: active-schema
	docker-compose stop postgres
	docker-compose rm -f postgres
	docker-compose up -d postgres

run-postgres:
	docker-compose stop postgres
	docker-compose rm -f postgres
	docker-compose up postgres

veryclean:

distribution:
	rm -rf distribution
	GOOS=darwin GOARCH=amd64 BUILD=$(BUILD)/distribution/darwin $(MAKE) $(BUILD)/distribution/darwin/fkflash
	cp -ar ~/.fk/tools $(BUILD)/distribution/darwin
	GOOS=windows GOARCH=amd64 BUILD=$(BUILD)/distribution/windows $(MAKE) $(BUILD)/distribution/windows/fkflash
	cp -ar ~/.fk/tools $(BUILD)/distribution/windows
