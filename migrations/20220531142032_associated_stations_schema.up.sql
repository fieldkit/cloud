CREATE TABLE fieldkit.associated_station (
    station_id INTEGER NOT NULL REFERENCES fieldkit.station (id),
    associated_station_id INTEGER NOT NULL REFERENCES fieldkit.station (id),
    priority INTEGER NOT NULL
);

CREATE UNIQUE INDEX associated_station_idx ON fieldkit.associated_station (station_id, associated_station_id);