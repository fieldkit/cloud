CREATE TABLE fieldkit.ttn_messages (
  id serial PRIMARY KEY,
  created_at timestamp NOT NULL DEFAULT NOW(),
  headers TEXT,
  body bytea NOT NULL
);
