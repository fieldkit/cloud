CREATE TABLE fieldkit.station (
  id integer NOT NULL,
  name text NOT NULL,
  user_id integer REFERENCES fieldkit.user (id) NOT NULL,
);

CREATE TABLE fieldkit.station_log (
  id integer NOT NULL,
  station_id integer NOT NULL,
  body text,
  timestamp TIMESTAMP NOT NULL
);
