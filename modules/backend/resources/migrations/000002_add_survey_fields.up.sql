-- Add created_by and tags fields to surveys table
ALTER TABLE surveys ADD COLUMN created_by VARCHAR(255);
ALTER TABLE surveys ADD COLUMN tags JSONB DEFAULT '[]'::jsonb;

-- Create index for tags to improve query performance
CREATE INDEX idx_surveys_tags ON surveys USING GIN (tags); 