modules = server/sqs-worker server/sqs-sender server/tools/fktool

GOARCH ?= amd64
GOOS ?= linux
GO ?= env GOOS=$(GOOS) GOARCH=$(GOARCH) go
BUILD ?= build

SERVER_SOURCES = $(shell find server -type f -name '*.go' -not -path "server/vendor/*")
TESTING_SOURCES = $(shell find testing -type f -name '*.go' -not -path "server/vendor/*")

default: binaries

binaries: $(BUILD)/server $(BUILD)/sqs-worker $(BUILD)/sqs-sender $(BUILD)/fktool $(BUILD)/fkflash $(BUILD)/testing-random $(BUILD)/weather-proxy $(BUILD)/inaturalist

all: binaries

tests:
	cd frontend && yarn run flow

$(BUILD)/server: $(SERVER_SOURCES) server/inaturalist/secrets.go
	$(GO) build -o $@ server/server.go

$(BUILD)/fktool: server/tools/fktool/*.go $(SERVER_SOURCES) $(TESTING_SOURCES)
	$(GO) build -o $@ server/tools/fktool/*.go

$(BUILD)/fkflash: server/tools/fkflash/*.go $(SERVER_SOURCES) $(TESTING_SOURCES)
	$(GO) build -o $@ server/tools/fkflash/*.go

$(BUILD)/inaturalist: server/tools/inaturalist/*.go $(SERVER_SOURCES) $(TESTING_SOURCES) server/inaturalist/secrets.go
	$(GO) build -o $@ server/tools/inaturalist/*.go

$(BUILD)/sqs-worker: server/sqs-worker/*.go $(SERVER_SOURCES)
	$(GO) build -o $@ server/sqs-worker/*.go

$(BUILD)/sqs-sender: server/sqs-sender/*.go $(SERVER_SOURCES)
	$(GO) build -o $@ server/sqs-sender/*.go

$(BUILD)/testing-random: testing/random/*.go $(SERVER_SOURCES) $(TESTING_SOURCES)
	$(GO) build -o $@ testing/random/*.go

$(BUILD)/weather-proxy: testing/weather-proxy/*.go $(SERVER_SOURCES) $(TESTING_SOURCES)
	$(GO) build -o $@ testing/weather-proxy/*.go

server/inaturalist/secrets.go: server/inaturalist/secrets.go.template
	cp server/inaturalist/secrets.go.template server/inaturalist/secrets.go

install: all
	cp build/fktool $(INSTALLDIR)
	cp build/testing-random $(INSTALLDIR)
	cp build/sqs-sender $(INSTALLDIR)
	cp build/sqs-worker $(INSTALLDIR)
	@for d in $(modules); do                           \
		(cd $$d && echo $$d && go install) || exit 1;  \
	done

generate:
	mv server/vendor server/vendor-temp # See https://github.com/goadesign/goa/issues/923
	(cd $(GOPATH)/src/github.com/fieldkit/cloud/server && go generate) || true
	mv server/vendor-temp server/vendor

deps: server/inaturalist/secrets.go
	cd server && go get ./...

clean:
	rm -rf build

clean-production:
	rm -rf schema-production

refresh-production: clean-production clone-production

schema-production:
	mkdir schema-production
	@if [ -d ~/conservify/dev-ops ]; then                                           \
		(cd ~/conservify/dev-ops/provisioning && ./db-dump.sh);                 \
		cp ~/conservify/dev-ops/schema.sql schema-production/000001.sql;        \
		cp ~/conservify/dev-ops/data.sql schema-production/000100.sql;          \
		cp schema/00000?.sql schema-production/;                                \
	else                                                                            \
		echo "No dev-ops directory found";                                      \
	fi

clone-production: schema-production
	rm -f active-schema
	ln -sf schema-production active-schema

restart-postgres:
	docker-compose stop postgres
	docker-compose rm -f postgres
	docker-compose up -d postgres

veryclean:

distribution:
	rm -rf distribution
	GOOS=darwin GOARCH=amd64 BUILD=build/distribution/darwin make build/distribution/darwin/fkflash
	cp -ar ~/.fk/tools build/distribution/darwin
	GOOS=windows GOARCH=amd64 BUILD=build/distribution/windows make build/distribution/windows/fkflash
	cp -ar ~/.fk/tools build/distribution/windows
