-- field notes

CREATE TABLE fieldkit.field_note_category (
    id serial PRIMARY KEY,
    key varchar(100) NOT NULL,
    name varchar(255) NOT NULL
);

CREATE TABLE fieldkit.field_note_media (
    id serial PRIMARY KEY,
    user_id integer REFERENCES fieldkit.user (id) NOT NULL,
    content_type varchar(100) NOT NULL,
    created timestamp NOT NULL,
    url varchar(255) NOT NULL
);

CREATE TABLE fieldkit.field_note (
    id serial PRIMARY KEY,
    created timestamp NOT NULL,
    station_id integer REFERENCES fieldkit.station (id) NOT NULL,
    user_id integer REFERENCES fieldkit.user (id) NOT NULL,
    category_id integer REFERENCES fieldkit.field_note_category (id) NOT NULL,
    media_id integer REFERENCES fieldkit.field_note_media (id),
    note text
);

INSERT INTO fieldkit.field_note_category (key, name) VALUES ('default', 'General');