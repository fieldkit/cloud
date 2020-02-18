FROM ubuntu AS source

WORKDIR /migrations

ADD https://github.com/golang-migrate/migrate/releases/download/v4.7.0/migrate.linux-amd64.tar.gz .

RUN tar xf migrate.linux-amd64.tar.gz && rm -rf *.tar.gz && mkdir bin && mv migrate.* bin

RUN ls -alhR

FROM scratch

COPY --from=source /migrations /migrations

COPY *.sql /migrations/

WORKDIR /migrations

ENTRYPOINT [ "/migrations/bin/migrate.linux-amd64" ]