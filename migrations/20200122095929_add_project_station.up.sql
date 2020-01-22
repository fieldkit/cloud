CREATE TABLE fieldkit.project_station (
	id serial PRIMARY KEY,
	station_id INTEGER REFERENCES fieldkit.station(id) NOT NULL,
	project_id INTEGER REFERENCES fieldkit.project(id) NOT NULL
);

CREATE UNIQUE INDEX ON fieldkit.project_station (station_id, project_id);
