-- ingestion

CREATE TABLE fieldkit.ingestion (
  id serial PRIMARY KEY,
  time timestamp NOT NULL,
  upload_id varchar(64) NOT NULL,
  user_id integer NOT NULL,
  device_id bytea NOT NULL,
  size integer NOT NULL,
  url varchar NOT NULL,
  blocks int8range NOT NULL,
  flags integer[] NOT NULL DEFAULT '{}'
);

CREATE INDEX ON fieldkit.ingestion (time, user_id);
