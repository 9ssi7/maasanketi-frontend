-- Drop index
DROP INDEX IF EXISTS idx_surveys_tags;

-- Remove columns
ALTER TABLE surveys DROP COLUMN IF EXISTS created_by;
ALTER TABLE surveys DROP COLUMN IF EXISTS tags; 