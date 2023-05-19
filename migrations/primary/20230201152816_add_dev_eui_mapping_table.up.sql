CREATE TABLE fieldkit.station_dev_eui (
        dev_eui bytea PRIMARY KEY,
        station_id integer NOT NULL REFERENCES fieldkit.station(id)
);