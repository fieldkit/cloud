CREATE TABLE fieldkit.station (
  id serial PRIMARY KEY,
  owner_id integer REFERENCES fieldkit.user (id) NOT NULL,
  device_id bytea NOT NULL,
  created_at timestamp NOT NULL DEFAULT NOW(),
  name text NOT NULL,
  status_json json NOT NULL
);

CREATE UNIQUE INDEX ON fieldkit.station (device_id);

CREATE TABLE fieldkit.station_log (
  id serial PRIMARY KEY,
  station_id integer NOT NULL REFERENCES fieldkit.station (id) NOT NULL,
  timestamp TIMESTAMP NOT NULL,
  body text
);
