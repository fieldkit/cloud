CREATE TABLE fieldkit.discussion_post (
    id serial PRIMARY KEY,
    user_id integer REFERENCES fieldkit.user (id) NOT NULL,
    thread_id integer REFERENCES fieldkit.discussion_post (id),
    project_id integer REFERENCES fieldkit.project (id),
	station_ids integer[] NOT NULL DEFAULT '{}',
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
	body text NOT NULL,
	context json
);
