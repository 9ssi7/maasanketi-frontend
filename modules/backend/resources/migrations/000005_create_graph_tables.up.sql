-- Create graphs table
CREATE TABLE IF NOT EXISTS graphs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    survey_id UUID NOT NULL REFERENCES surveys(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    kind VARCHAR(50) NOT NULL,
    content JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Create response_graphs table
CREATE TABLE IF NOT EXISTS response_graphs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    survey_id UUID NOT NULL REFERENCES surveys(id) ON DELETE CASCADE,
    graph_id UUID NOT NULL REFERENCES graphs(id) ON DELETE CASCADE,
    content JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Create indexes
CREATE INDEX idx_graphs_survey_id ON graphs(survey_id);
CREATE INDEX idx_response_graphs_survey_id ON response_graphs(survey_id);
CREATE INDEX idx_response_graphs_graph_id ON response_graphs(graph_id);
CREATE INDEX idx_graphs_created_at ON graphs(created_at);
CREATE INDEX idx_response_graphs_created_at ON response_graphs(created_at); 