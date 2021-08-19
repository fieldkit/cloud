CREATE TABLE fieldkit.station_model (
  id serial PRIMARY KEY,
  name TEXT NOT NULL
);

INSERT INTO fieldkit.station_model (name) VALUES ('FieldKit Station');
INSERT INTO fieldkit.station_model (name) VALUES ('FloodNet');

ALTER TABLE fieldkit.station ADD COLUMN model_id INTEGER REFERENCES fieldkit.station_model(id);
UPDATE fieldkit.station SET model_id = 1;
ALTER TABLE fieldkit.station ALTER model_id SET NOT NULL;