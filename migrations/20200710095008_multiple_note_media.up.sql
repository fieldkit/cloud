/* This table is empty, unused so far. */
ALTER TABLE fieldkit.notes DROP COLUMN media_id;

CREATE TABLE fieldkit.notes_media_link (
	note_id integer REFERENCES fieldkit.notes (id),
    media_id integer REFERENCES fieldkit.notes_media (id)
);

CREATE UNIQUE INDEX ON fieldkit.notes_media_link (note_id, media_id);
