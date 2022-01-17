CREATE TABLE fieldkit.notification (
  id serial PRIMARY KEY,
  created_at timestamp NOT NULL DEFAULT NOW(),
  user_id integer REFERENCES fieldkit.user(id) NOT NULL,
  post_id integer REFERENCES fieldkit.discussion_post(id),
  key text NOT NULL,
  kind text NOT NULL,
  body bytea,
  seen BOOL NOT NULL
);
