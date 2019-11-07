CREATE TABLE fieldkit.project_invite (
    id serial PRIMARY KEY,
    project_id integer REFERENCES fieldkit.project (id) NOT NULL,
    user_id integer REFERENCES fieldkit.user (id) NOT NULL,
    invited_email varchar(255) NOT NULL,
    invited_time timestamp NOT NULL,
    accepted_time timestamp,
    rejected_time timestamp
);