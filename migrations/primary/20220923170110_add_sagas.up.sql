CREATE TABLE fieldkit.sagas (
    id TEXT NOT NULL PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
    version INTEGER NOT NULL,
	scheduled_at TIMESTAMP,
	tags JSONB NOT NULL,
    type TEXT NOT NULL,
	body JSONB NOT NULL
);
