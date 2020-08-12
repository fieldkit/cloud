DROP TABLE fieldkit.data_export;

CREATE TABLE fieldkit.data_export (
    id serial PRIMARY KEY,
	token bytea NOT NULL,
    user_id INTEGER REFERENCES fieldkit.user (id) NOT NULL,
    kind TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
	completed_at TIMESTAMP,
	progress FLOAT NOT NULL,
	download_url TEXT,
    args JSON NOT NULL
);

CREATE UNIQUE INDEX ON fieldkit.data_export(token);
