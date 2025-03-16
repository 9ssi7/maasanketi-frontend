-- Drop indexes
DROP INDEX IF EXISTS idx_survey_responses_survey_id;
DROP INDEX IF EXISTS idx_survey_responses_user_id;
DROP INDEX IF EXISTS idx_survey_responses_is_completed;
DROP INDEX IF EXISTS idx_surveys_created_at;
DROP INDEX IF EXISTS idx_survey_responses_created_at;

-- Drop tables
DROP TABLE IF EXISTS survey_responses;
DROP TABLE IF EXISTS surveys;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp"; 