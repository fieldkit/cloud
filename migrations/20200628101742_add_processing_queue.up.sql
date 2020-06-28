CREATE TABLE fieldkit.ingestion_queue (
	id SERIAL PRIMARY KEY,
	ingestion_id integer NOT NULL REFERENCES fieldkit.ingestion(id),
	queued TIMESTAMP NOT NULL,
	attempted TIMESTAMP,
	completed TIMESTAMP,
	total_records INTEGER,
	other_errors INTEGER,
	meta_errors INTEGER,
	data_errors INTEGER
);

CREATE INDEX ON fieldkit.ingestion_queue(queued, completed);
