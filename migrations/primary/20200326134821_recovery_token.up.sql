CREATE TABLE fieldkit.recovery_token (
	token bytea PRIMARY KEY,
	user_id integer REFERENCES fieldkit.user (id) ON DELETE CASCADE NOT NULL,
	expires timestamp NOT NULL
);

CREATE UNIQUE INDEX ON fieldkit.recovery_token (user_id);
