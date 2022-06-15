ALTER TABLE fieldkit.ttn_schema ADD COLUMN project_id INTEGER REFERENCES fieldkit.project (id);

UPDATE fieldkit.ttn_schema SET project_id = 174 WHERE id IN (7, 8);