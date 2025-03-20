package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/mstrYoda/maasanketi.co/domain/graph"
	"github.com/mstrYoda/maasanketi.co/domain/resgraph"
	"github.com/mstrYoda/maasanketi.co/domain/response"
	"github.com/mstrYoda/maasanketi.co/domain/survey"
	log "github.com/mstrYoda/maasanketi.co/pkg/logger"
	"github.com/mstrYoda/maasanketi.co/repository"
	"go.uber.org/zap"
)

// GraphGenerationService handles generating graphs from survey responses
type GraphGenerationService struct {
	repo *repository.Repository
}

// NewGraphGenerationService creates a new GraphGenerationService
func NewGraphGenerationService(repo *repository.Repository) *GraphGenerationService {
	return &GraphGenerationService{
		repo: repo,
	}
}

// GenerateGraphsForCompletedSurveys finds all surveys that have ended and generates graphs for them
func (s *GraphGenerationService) GenerateGraphsForCompletedSurveys(ctx context.Context) error {
	page := uint64(0)
	limit := uint64(10)

	for {
		surveys, err := s.repo.Survey.FindCompletedSurveys(ctx, page, limit)
		if err != nil {
			return fmt.Errorf("error finding completed surveys: %w", err)
		}

		if len(surveys) == 0 {
			log.Logger().Info("no completed surveys found for graph generation")
			break
		}

		log.Logger().Info("found completed surveys for graph generation", zap.Int("count", len(surveys)))

		for _, survey := range surveys {
			if err := s.generateGraphsForSurvey(ctx, survey); err != nil {
				log.Logger().Error("error generating graphs for survey",
					zap.String("survey_id", survey.ID),
					zap.Error(err))
				continue
			}
		}

		page++
	}

	return nil
}

// generateGraphsForSurvey generates all graphs for a single survey
func (s *GraphGenerationService) generateGraphsForSurvey(ctx context.Context, survey *survey.Survey) error {
	// Get all graphs defined for this survey
	graphs, err := s.repo.Graph.ListBySurveyID(ctx, survey.ID)
	if err != nil {
		return fmt.Errorf("failed to get graphs for survey %s: %w", survey.ID, err)
	}

	// Get all responses for this survey
	responses, err := s.repo.Response.List(ctx, survey.ID)
	if err != nil {
		return fmt.Errorf("failed to get responses for survey %s: %w", survey.ID, err)
	}

	if len(responses) == 0 {
		log.Logger().Warn("no responses found for survey, skipping graph generation",
			zap.String("survey_id", survey.ID))
		return nil
	}

	// Process each graph
	for _, g := range graphs {
		if err := s.processGraph(ctx, g, responses); err != nil {
			log.Logger().Error("error processing graph",
				zap.String("graph_id", g.ID),
				zap.Error(err))
			continue
		}
	}

	return nil
}

// processGraph processes a single graph for a survey
func (s *GraphGenerationService) processGraph(ctx context.Context, g *graph.Graph, responses []*response.Response) error {
	// Create a new response graph
	responseGraph := &resgraph.ResponseGraph{
		ID:        uuid.New().String(),
		SurveyID:  g.SurveyID,
		GraphID:   g.ID,
		Content:   []resgraph.ResponseGraphContent{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Process the graph based on its type and strategy
	content, err := s.generateGraphContent(g, responses)
	if err != nil {
		return fmt.Errorf("failed to generate graph content: %w", err)
	}

	responseGraph.Content = content

	// Save the response graph
	if err := s.repo.ResponseGraph.Create(ctx, responseGraph); err != nil {
		return fmt.Errorf("failed to save response graph: %w", err)
	}

	log.Logger().Info("generated graph for survey",
		zap.String("graph_id", g.ID),
		zap.String("survey_id", g.SurveyID))
	return nil
}

// generateGraphContent creates the content for a graph based on its configuration
func (s *GraphGenerationService) generateGraphContent(g *graph.Graph, responses []*response.Response) ([]resgraph.ResponseGraphContent, error) {
	// Map to store aggregated data
	data := make(map[string]map[string][]interface{})

	// Initialize data structure
	for _, key := range g.Content.KeyFields {
		data[key.Field] = make(map[string][]interface{})
	}

	// Process responses and collect data
	for _, resp := range responses {
		if !resp.IsCompleted {
			continue
		}

		// For each key field in the graph
		for _, key := range g.Content.KeyFields {
			keyValue, ok := resp.Answers[key.Field]
			if !ok {
				continue
			}

			keyValueStr := fmt.Sprintf("%v", keyValue)

			// For each value field in the graph
			for _, value := range g.Content.ValueFields {
				val, ok := resp.Answers[value.Field]
				if !ok {
					continue
				}

				// Store the value for later aggregation
				if _, exists := data[key.Field][keyValueStr]; !exists {
					data[key.Field][keyValueStr] = make([]interface{}, 0)
				}

				data[key.Field][keyValueStr] = append(data[key.Field][keyValueStr], val)
			}
		}
	}

	// Generate content for the graph
	var content []resgraph.ResponseGraphContent

	// For each key field
	for _, key := range g.Content.KeyFields {
		var labels []string
		var values []string

		// Sort keys for consistent output
		keyData := data[key.Field]

		for keyStr, vals := range keyData {
			labels = append(labels, keyStr)

			// Process each value field
			for i, value := range g.Content.ValueFields {
				if i >= len(vals) {
					continue
				}

				// Apply the appropriate strategy to the values
				aggregated := s.applyStrategy(value.Strategy, vals)
				values = append(values, fmt.Sprintf("%v", aggregated))
			}
		}

		if len(labels) == 0 || len(values) == 0 {
			log.Logger().Warn("no labels or values found for graph",
				zap.String("graph_id", g.ID),
				zap.String("survey_id", g.SurveyID))
			continue
		}

		content = append(content, resgraph.ResponseGraphContent{
			Labels: labels,
			Values: values,
		})
	}

	return content, nil
}

// applyStrategy applies the specified aggregation strategy to a set of values
func (s *GraphGenerationService) applyStrategy(strategy graph.ValueStrategy, values []interface{}) interface{} {
	switch strategy {
	case graph.ValueStrategyCount:
		return len(values)

	case graph.ValueStrategySum:
		var sum float64
		for _, v := range values {
			// Try to convert to float
			f, err := toFloat(v)
			if err == nil {
				sum += f
			}
		}
		return sum

	case graph.ValueStrategyAvg:
		var sum float64
		count := 0
		for _, v := range values {
			// Try to convert to float
			f, err := toFloat(v)
			if err == nil {
				sum += f
				count++
			}
		}
		if count > 0 {
			return sum / float64(count)
		}
		return 0

	default:
		return len(values)
	}
}

// toFloat tries to convert an interface{} to a float64
func toFloat(v interface{}) (float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case float32:
		return float64(val), nil
	case int:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	case json.Number:
		return val.Float64()
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", v)
	}
}
