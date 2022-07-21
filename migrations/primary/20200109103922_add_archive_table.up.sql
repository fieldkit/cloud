CREATE TABLE fieldkit.archive_history (
	id serial PRIMARY KEY,
	archived timestamp NOT NULL DEFAULT now(),
	old_device_id bytea NOT NULL,
	new_device_id bytea NOT NULL
);
