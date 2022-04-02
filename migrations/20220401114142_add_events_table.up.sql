CREATE TABLE fieldkit.data_event (
    id serial PRIMARY KEY,
    user_id integer REFERENCES fieldkit.user (id) NOT NULL,
    project_id integer REFERENCES fieldkit.project (id),
	station_ids integer[] NOT NULL DEFAULT '{}',
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    start_time timestamp NOT NULL,
    end_time timestamp NOT NULL,
    title text NOT NULL,
	description text NOT NULL,
	context json
);

CREATE INDEX ON fieldkit.data_event (project_id, start_time, end_time)