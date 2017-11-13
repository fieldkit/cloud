CREATE TABLE fieldkit.raw_message (
    id serial PRIMARY KEY,
    time timestamp,
    origin_id varchar NOT NULL, -- SQS Id or some other provider generated identifier.
    data json NOT NULL
);

CREATE TABLE fieldkit.device (
    input_id integer REFERENCES fieldkit.input (id) ON DELETE CASCADE PRIMARY KEY,
    key varchar NOT NULL,
    token bytea NOT NULL
);

CREATE INDEX ON fieldkit.device (key);

CREATE INDEX ON fieldkit.device (token);

CREATE TABLE fieldkit.device_schema (
    id serial PRIMARY KEY,
    device_id integer REFERENCES fieldkit.device (input_id) ON DELETE CASCADE,
    schema_id integer REFERENCES fieldkit.schema (id) ON DELETE CASCADE,
    key varchar
);

CREATE UNIQUE INDEX ON fieldkit.device_schema (device_id, schema_id);

CREATE TABLE fieldkit.device_location (
  id serial PRIMARY KEY,
  timestamp timestamp NOT NULL,
  device_id integer REFERENCES fieldkit.device (input_id) NOT NULL,
  location geometry(POINT, 4326) NOT NULL
);

CREATE INDEX ON fieldkit.device_location (device_id, timestamp);
