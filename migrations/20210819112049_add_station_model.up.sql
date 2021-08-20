CREATE TABLE fieldkit.ttn_schema (
  id serial PRIMARY KEY,
  owner_id INTEGER NOT NULL REFERENCES fieldkit.user(id),
  token BYTEA NOT NULL,
  name TEXT NOT NULL,
  body JSON NOT NULL 
);

ALTER TABLE fieldkit.ttn_messages ADD COLUMN schema_id INTEGER REFERENCES fieldkit.ttn_schema(id);

CREATE TABLE fieldkit.station_model (
  id serial PRIMARY KEY,
  ttn_schema_id INTEGER REFERENCES fieldkit.ttn_schema(id),
  name TEXT NOT NULL
);

INSERT INTO fieldkit.station_model (name, ttn_schema_id) VALUES ('FieldKit Station', NULL);

ALTER TABLE fieldkit.station ADD COLUMN model_id INTEGER REFERENCES fieldkit.station_model(id);
UPDATE fieldkit.station SET model_id = 1;
ALTER TABLE fieldkit.station ALTER model_id SET NOT NULL;