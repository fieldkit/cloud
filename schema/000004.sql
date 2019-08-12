CREATE TABLE fieldkit.station (
  id integer NOT NULL,
  name integer NOT NULL
);

CREATE TABLE fieldkit.station_log (
  id integer NOT NULL,
  station_id integer NOT NULL,
  body text,
  timestamp TIMESTAMP NOT NULL
);
