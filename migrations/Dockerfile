FROM ubuntu AS source

WORKDIR /migrations

ADD https://github.com/golang-migrate/migrate/releases/download/v4.7.0/migrate.linux-amd64.tar.gz /migrations

RUN tar xf migrate.linux-amd64.tar.gz && rm -rf *.tar.gz && mv migrate.* bin

RUN ls -alh

FROM scratch

COPY --from=source /migrations /migrations
COPY *.sql /migrations/

WORKDIR /migrations

ENTRYPOINT /migrate/bin/migrate.linux-amd64 --path /migrations --database "$DATABASE_URL"