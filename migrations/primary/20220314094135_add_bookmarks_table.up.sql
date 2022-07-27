CREATE TABLE fieldkit.bookmarks (
	id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES fieldkit.user(id),
    token TEXT NOT NULL,
    bookmark TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    referenced_at TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX ON fieldkit.bookmarks (token);