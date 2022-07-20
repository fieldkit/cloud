CREATE TABLE fieldkit.data_export (
    id serial PRIMARY KEY,
	token bytea NOT NULL,
    user_id INTEGER REFERENCES fieldkit.user (id) NOT NULL,
    created_at TIMESTAMP NOT NULL,
	completed_at TIMESTAMP,
	download_url TEXT,
	progress FLOAT NOT NULL,
    args JSON NOT NULL
);

CREATE UNIQUE INDEX ON fieldkit.data_export(token);
