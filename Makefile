modules = server/sqs-worker server/sqs-sender server/api/tool/fkcli server/api/tool/fktool

BUILD=build

SERVER_SOURCES = $(shell find server -type f -name '*.go' -not -path "server/vendor/*")
TESTING_SOURCES = $(shell find testing -type f -name '*.go' -not -path "server/vendor/*")

default: binaries tests

binaries: $(BUILD)/server $(BUILD)/sqs-worker $(BUILD)/sqs-sender $(BUILD)/fkcli $(BUILD)/fktool $(BUILD)/testing-random $(BUILD)/weather-proxy $(BUILD)/naturalist

all: binaries

tests:
	cd frontend && yarn run flow

$(BUILD)/server: $(SERVER_SOURCES)
	go build -o $@ server/server.go

$(BUILD)/sqs-worker: server/sqs-worker/*.go $(SERVER_SOURCES)
	go build -o $@ server/sqs-worker/*.go

$(BUILD)/sqs-sender: server/sqs-sender/*.go $(SERVER_SOURCES)
	go build -o $@ server/sqs-sender/*.go

$(BUILD)/fkcli: server/api/tool/fkcli/*.go $(SERVER_SOURCES)
	go build -o $@ server/api/tool/fkcli/*.go

$(BUILD)/fktool: server/api/tool/fktool/*.go $(SERVER_SOURCES) $(TESTING_SOURCES)
	go build -o $@ server/api/tool/fktool/*.go

$(BUILD)/testing-random: testing/random/*.go $(SERVER_SOURCES) $(TESTING_SOURCES)
	go build -o $@ testing/random/*.go

$(BUILD)/weather-proxy: testing/weather-proxy/*.go $(SERVER_SOURCES) $(TESTING_SOURCES)
	go build -o $@ testing/weather-proxy/*.go

$(BUILD)/naturalist: server/tools/naturalist/*.go $(SERVER_SOURCES) $(TESTING_SOURCES)
	go build -o $@ server/tools/naturalist/*.go

server/naturalist/secrets.go: server/naturalist/secrets.go.template
	cp server/naturalist/secrets.go.template server/naturalist/secrets.go

install: all
	cp build/fktool $(INSTALLDIR)
	cp build/fkcli $(INSTALLDIR)
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

deps:
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
