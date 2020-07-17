ALTER TABLE fieldkit.notes ADD COLUMN version INTEGER NOT NULL DEFAULT 0;
ALTER TABLE fieldkit.notes ADD COLUMN updated_at TIMESTAMP;
UPDATE fieldkit.notes SET updated_at = created_at;
ALTER TABLE fieldkit.notes ALTER COLUMN updated_at SET NOT NULL;
