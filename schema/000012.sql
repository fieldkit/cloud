CREATE TABLE fieldkit.station (
  id integer NOT NULL,
  owner_id integer REFERENCES fieldkit.user (id) NOT NULL,
  device_id bytea NOT NULL,
  created_at timestamp NOT NULL DEFAULT NOW(),
  name text NOT NULL
);

CREATE UNIQUE INDEX ON fieldkit.station (device_id);

CREATE TABLE fieldkit.station_log (
  id integer NOT NULL,
  station_id integer NOT NULL,
  body text,
  timestamp TIMESTAMP NOT NULL
);
