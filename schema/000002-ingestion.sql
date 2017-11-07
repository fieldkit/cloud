CREATE TABLE fieldkit.raw_message (
    id serial PRIMARY KEY,
    time timestamp,
    origin_id varchar NOT NULL, -- SQS Id or some other provider generated identifier.
    data json NOT NULL,
    hash varchar NOT NULL
);

CREATE TABLE fieldkit.device (
    input_id integer REFERENCES fieldkit.input (id) ON DELETE CASCADE PRIMARY KEY,
    token bytea NOT NULL
);

CREATE INDEX ON fieldkit.device (token);

CREATE TABLE fieldkit.stream (
    id serial PRIMARY KEY,
    name varchar(256) NOT NULL,
    device_id integer REFERENCES fieldkit.device (input_id) NOT NULL,
    schema_id integer REFERENCES fieldkit.schema (id) ON DELETE CASCADE NOT NULL
);

CREATE INDEX ON fieldkit.stream (device_id);
