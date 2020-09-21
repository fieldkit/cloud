CREATE TABLE fieldkit.project_station_new (
	project_id INTEGER REFERENCES fieldkit.project(id) NOT NULL,
	station_id INTEGER REFERENCES fieldkit.station(id) NOT NULL,
	PRIMARY KEY (project_id, station_id)
);

INSERT INTO fieldkit.project_station_new (project_id, station_id) SELECT project_id, station_id FROM fieldkit.project_station;

ALTER TABLE fieldkit.project_station RENAME TO project_station_old;

ALTER TABLE fieldkit.project_station_new RENAME TO project_station;
