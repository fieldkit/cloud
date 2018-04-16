CREATE TABLE fieldkit.inaturalist_observations (
  id integer PRIMARY KEY NOT NULL,
  updated_at timestamp NOT NULL,
	timestamp timestamp NOT NULL,
	location geometry(POINT, 4326) NOT NULL,
	data jsonb NOT NULL
);

CREATE INDEX ON fieldkit.inaturalist_observations (timestamp);
CREATE INDEX ON fieldkit.inaturalist_observations USING GIST (location);
