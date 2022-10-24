CREATE TABLE fieldkit.sensor_data_dirty (
    id SERIAL PRIMARY KEY,
    modified TIMESTAMPTZ NOT NULL,
    data_start TIMESTAMPTZ NOT NULL,
    data_end TIMESTAMPTZ NOT NULL
);