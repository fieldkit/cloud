/* This table is empty, unused so far. */
ALTER TABLE fieldkit.notes DROP COLUMN media_id;

CREATE TABLE fieldkit.notes_media_link (
	note_id integer NOT NULL REFERENCES fieldkit.notes (id),
    media_id integer NOT NULL REFERENCES fieldkit.notes_media (id)
);

CREATE UNIQUE INDEX ON fieldkit.notes_media_link (note_id, media_id);

/* This table is empty, unused so far. */
ALTER TABLE fieldkit.notes_media ADD COLUMN key TEXT NOT NULL;
ALTER TABLE fieldkit.notes_media ADD COLUMN station_id INTEGER NOT NULL REFERENCES fieldkit.station (id);

CREATE UNIQUE INDEX ON fieldkit.notes_media (station_id, key);
