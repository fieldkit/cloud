BEGIN;
-- Protecting from inserts that may happen while we update, though for these
-- tables that should never happen. Keep this here to maintain the correct
-- pattern, though.
LOCK TABLE fieldkit.sensor_meta IN EXCLUSIVE MODE;
SELECT setval('sensor_meta_id_seq', COALESCE((SELECT MAX(id) + 1 FROM fieldkit.sensor_meta), 1), false);

LOCK TABLE fieldkit.module_meta IN EXCLUSIVE MODE;
SELECT setval('module_meta_id_seq', COALESCE((SELECT MAX(id) + 1 FROM fieldkit.module_meta), 1), false);
COMMIT;