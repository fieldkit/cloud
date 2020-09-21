CREATE TABLE fieldkit.station_owner (
	id serial PRIMARY KEY,
	station_id integer REFERENCES fieldkit.station (id) NOT NULL,
	user_id integer REFERENCES fieldkit.user (id) NOT NULL,
	start_time timestamp NOT NULL,
	end_time timestamp NOT NULL
);
