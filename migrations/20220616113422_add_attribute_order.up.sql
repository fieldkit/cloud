ALTER TABLE fieldkit.project_attribute ADD COLUMN priority INTEGER;
UPDATE fieldkit.project_attribute SET priority = 10;
UPDATE fieldkit.project_attribute SET priority = 0 WHERE name = 'Neighborhood';
UPDATE fieldkit.project_attribute SET priority = 1 WHERE name = 'Borough';
ALTER TABLE fieldkit.project_attribute ALTER priority SET NOT NULL;