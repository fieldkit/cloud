CREATE TABLE fieldkit.inaturalist_observations (
  id integer NOT NULL,
  site_id integer NOT NULL,
  updated_at timestamp NOT NULL,
	timestamp timestamp NOT NULL,
	location geometry(POINT, 4326) NOT NULL,
	data jsonb NOT NULL,
  PRIMARY KEY (id, site_id)
);

CREATE INDEX IF NOT EXISTS inaturalist_observations_timestamp_idx ON fieldkit.inaturalist_observations (timestamp);
CREATE INDEX IF NOT EXISTS inaturalist_observations_location_idx ON fieldkit.inaturalist_observations USING GIST (location);
