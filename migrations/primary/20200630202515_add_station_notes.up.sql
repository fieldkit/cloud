CREATE TABLE fieldkit.notes_media (
    id serial PRIMARY KEY,
    user_id integer REFERENCES fieldkit.user (id) NOT NULL,
    content_type varchar(100) NOT NULL,
    created timestamp NOT NULL,
    url varchar(255) NOT NULL
);

CREATE TABLE fieldkit.notes (
    id serial PRIMARY KEY,
    created_at timestamp NOT NULL,
    station_id integer REFERENCES fieldkit.station (id) NOT NULL,
    author_id integer REFERENCES fieldkit.user (id) NOT NULL,
    media_id integer REFERENCES fieldkit.notes_media (id),
	key text,
    body text
);

ALTER TABLE fieldkit.station ADD COLUMN photo_id INTEGER REFERENCES fieldkit.notes_media (id);
