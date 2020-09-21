CREATE TABLE fieldkit.data_visibility (
	id serial PRIMARY KEY,
	start_time timestamp NOT NULL,
	end_time timestamp NOT NULL,
	station_id integer REFERENCES fieldkit.station (id) NOT NULL,
	project_id integer REFERENCES fieldkit.project (id),
	user_id integer REFERENCES fieldkit.user (id)
);

CREATE UNIQUE INDEX ON fieldkit.data_visibility (start_time, end_time, station_id);
