-- Drop indexes
DROP INDEX IF EXISTS idx_response_graphs_created_at;
DROP INDEX IF EXISTS idx_graphs_created_at;
DROP INDEX IF EXISTS idx_response_graphs_graph_id;
DROP INDEX IF EXISTS idx_response_graphs_survey_id;
DROP INDEX IF EXISTS idx_graphs_survey_id;

-- Drop tables
DROP TABLE IF EXISTS response_graphs;
DROP TABLE IF EXISTS graphs; 