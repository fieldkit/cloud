CREATE TABLE fieldkit.collection (
	id serial PRIMARY KEY,
	owner_id integer REFERENCES fieldkit.user (id) NOT NULL,
	created_at timestamp NOT NULL DEFAULT NOW(),
	updated_at timestamp NOT NULL DEFAULT NOW(),
	name text NOT NULL,
	description text NOT NULL DEFAULT '',
	origin text NOT NULL,
	tags text NOT NULL DEFAULT '',
	private boolean NOT NULL DEFAULT false
);

CREATE TABLE fieldkit.collection_viewer (
	collection_id integer REFERENCES fieldkit.collection (id) NOT NULL,
	user_id integer REFERENCES fieldkit.user (id) NOT NULL,
	PRIMARY KEY (collection_id, user_id)
);

CREATE TABLE fieldkit.collection_station (
	collection_id integer REFERENCES fieldkit.collection (id) NOT NULL,
	station_id INTEGER REFERENCES fieldkit.station(id) NOT NULL,
	PRIMARY KEY (collection_id, station_id)
);
