FROM golang AS source

WORKDIR /work

ADD go.* /work/
ADD *.go /work/

RUN ls -alhR

RUN go get ./...
RUN go build -o migrate *.go

FROM ubuntu

COPY --from=source /work/migrate /migrations/bin/

COPY *.sql /migrations/

WORKDIR /migrations

ENTRYPOINT [ "/migrations/bin/migrate" ]