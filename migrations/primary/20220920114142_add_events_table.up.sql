CREATE TABLE fieldkit.data_event (
    id serial PRIMARY KEY,
    user_id integer REFERENCES fieldkit.user (id) NOT NULL,
    project_ids integer[] NOT NULL DEFAULT '{}',
	station_ids integer[] NOT NULL DEFAULT '{}',
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    start_time timestamp NOT NULL,
    end_time timestamp NOT NULL,
    title text NOT NULL,
	description text NOT NULL,
	context json
);

CREATE INDEX ON fieldkit.data_event (start_time, end_time)