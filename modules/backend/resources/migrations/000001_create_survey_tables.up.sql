-- Create extension for UUID support
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create surveys table
CREATE TABLE IF NOT EXISTS surveys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    slug VARCHAR(255) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    questions JSONB NOT NULL,
    min_completion_time_min INTEGER NOT NULL DEFAULT 3,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Create survey_responses table
CREATE TABLE IF NOT EXISTS survey_responses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    survey_id UUID NOT NULL REFERENCES surveys(id) ON DELETE CASCADE,
    user_id UUID,
    answers JSONB NOT NULL,
    started_at TIMESTAMP WITH TIME ZONE NOT NULL,
    completed_at TIMESTAMP WITH TIME ZONE,
    ip_address VARCHAR(45) NOT NULL,
    user_agent TEXT,
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    is_anonymous BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Create indexes
CREATE INDEX idx_survey_responses_survey_id ON survey_responses(survey_id);
CREATE INDEX idx_survey_responses_user_id ON survey_responses(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_survey_responses_is_completed ON survey_responses(is_completed);
CREATE INDEX idx_surveys_created_at ON surveys(created_at);
CREATE INDEX idx_survey_responses_created_at ON survey_responses(created_at);
CREATE INDEX idx_surveys_slug ON surveys(slug); 