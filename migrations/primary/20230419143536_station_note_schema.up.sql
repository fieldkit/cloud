CREATE TABLE fieldkit.station_note (
    id serial PRIMARY KEY,
    station_id integer REFERENCES fieldkit.station (id) NOT NULL,
    user_id integer REFERENCES fieldkit.user (id) NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
	body text NOT NULL
);
