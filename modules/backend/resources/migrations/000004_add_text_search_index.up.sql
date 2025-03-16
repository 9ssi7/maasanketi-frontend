-- Add a GIN index for text search on title and description
CREATE INDEX idx_surveys_text_search ON surveys USING GIN (to_tsvector('simple', title || ' ' || description)); 