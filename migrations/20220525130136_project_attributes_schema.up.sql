CREATE TABLE fieldkit.project_attribute (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES fieldkit.project (id),
    name TEXT NOT NULL
);

CREATE UNIQUE INDEX project_attribute_idx ON fieldkit.project_attribute (project_id, name);

CREATE TABLE fieldkit.station_project_attribute (
    id SERIAL PRIMARY KEY,
    station_id INTEGER NOT NULL REFERENCES fieldkit.station (id),
    attribute_id INTEGER NOT NULL REFERENCES fieldkit.project_attribute (id),
    string_value TEXT NOT NULL
);

CREATE UNIQUE INDEX station_project_attribute_idx ON fieldkit.station_project_attribute (station_id, attribute_id);