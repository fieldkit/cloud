CREATE TABLE fieldkit.module_meta (
	id SERIAL PRIMARY KEY,
    key TEXT NOT NULL,
    manufacturer INTEGER NOT NULL,
    kinds INTEGER[] NOT NULL,
    version INTEGER[] NOT NULL,
    internal BOOL NOT NULL
);

CREATE TABLE fieldkit.sensor_meta (
	id SERIAL PRIMARY KEY,
	module_id INTEGER NOT NULL REFERENCES fieldkit.module_meta(id),
    ordering INTEGER NOT NULL,
    sensor_key TEXT NOT NULL,
    firmware_key TEXT NOT NULL,
    full_key TEXT NOT NULL,
    internal BOOL NOT NULL,
    uom TEXT NOT NULL,
    strings JSON NOT NULL,
    viz JSON NOT NULL,
    ranges JSON NOT NULL
);