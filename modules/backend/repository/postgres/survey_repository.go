package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mstrYoda/maasanketi.co/entity"
	"github.com/mstrYoda/maasanketi.co/pkg/list"
	"github.com/mstrYoda/maasanketi.co/pkg/rescode"
)

// SurveyRepository is a PostgreSQL implementation of the survey repository
type SurveyRepository struct {
	db *pgxpool.Pool
	sb squirrel.StatementBuilderType
}

// NewSurveyRepository creates a new PostgreSQL survey repository
func NewSurveyRepository(db *pgxpool.Pool) *SurveyRepository {
	return &SurveyRepository{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// parseUUID parses a string to UUID
func parseUUID(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}

// CreateSurvey creates a new survey
func (r *SurveyRepository) CreateSurvey(ctx context.Context, survey *entity.Survey) error {
	questionsJSON, err := json.Marshal(survey.Questions)
	if err != nil {
		return rescode.Failed(err)
	}

	tagsJSON, err := json.Marshal(survey.Tags)
	if err != nil {
		return rescode.Failed(err)
	}

	// Parse UUID from string
	surveyUUID, err := parseUUID(survey.ID)
	if err != nil {
		return rescode.IDInvalid(err)
	}

	// Generate slug if not provided
	if survey.Slug == "" {
		survey.GenerateSlug()
	}

	query := r.sb.Insert("surveys").
		Columns("id", "slug", "title", "description", "questions", "min_completion_time_min", "created_at", "updated_at", "created_by", "tags", "finishes_at").
		Values(surveyUUID, survey.Slug, survey.Title, survey.Description, questionsJSON, survey.MinCompletionTimeMin, survey.CreatedAt, survey.UpdatedAt, survey.CreatedBy, tagsJSON, survey.FinishesAt)

	sql, args, err := query.ToSql()
	if err != nil {
		return rescode.Failed(err)
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return rescode.Failed(err)
	}

	return nil
}

// GetSurveyByID gets a survey by ID
func (r *SurveyRepository) GetSurveyByID(ctx context.Context, id string) (*entity.Survey, error) {
	// Parse UUID from string
	surveyUUID, err := parseUUID(id)
	if err != nil {
		return nil, rescode.IDInvalid(err)
	}

	query := r.sb.Select("id", "slug", "title", "description", "questions", "min_completion_time_min", "created_at", "updated_at", "created_by", "tags", "finishes_at").
		From("surveys").
		Where(squirrel.Eq{"id": surveyUUID})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, rescode.Failed(err)
	}

	var survey entity.Survey
	var questionsJSON []byte
	var tagsJSON []byte
	var dbUUID uuid.UUID

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&dbUUID,
		&survey.Slug,
		&survey.Title,
		&survey.Description,
		&questionsJSON,
		&survey.MinCompletionTimeMin,
		&survey.CreatedAt,
		&survey.UpdatedAt,
		&survey.CreatedBy,
		&tagsJSON,
		&survey.FinishesAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, rescode.SurveyNotFound(err)
		}
		return nil, rescode.Failed(err)
	}

	survey.ID = dbUUID.String()

	err = json.Unmarshal(questionsJSON, &survey.Questions)
	if err != nil {
		return nil, rescode.Failed(err)
	}

	err = json.Unmarshal(tagsJSON, &survey.Tags)
	if err != nil {
		return nil, rescode.Failed(err)
	}

	return &survey, nil
}

// UpdateSurvey updates a survey
func (r *SurveyRepository) UpdateSurvey(ctx context.Context, survey *entity.Survey) error {
	questionsJSON, err := json.Marshal(survey.Questions)
	if err != nil {
		return rescode.Failed(err)
	}

	tagsJSON, err := json.Marshal(survey.Tags)
	if err != nil {
		return rescode.Failed(err)
	}

	// Parse UUID from string
	surveyUUID, err := parseUUID(survey.ID)
	if err != nil {
		return rescode.IDInvalid(err)
	}

	query := r.sb.Update("surveys").
		Set("title", survey.Title).
		Set("description", survey.Description).
		Set("questions", questionsJSON).
		Set("min_completion_time_min", survey.MinCompletionTimeMin).
		Set("updated_at", survey.UpdatedAt).
		Set("created_by", survey.CreatedBy).
		Set("tags", tagsJSON).
		Set("finishes_at", survey.FinishesAt).
		Where(squirrel.Eq{"id": surveyUUID})

	sql, args, err := query.ToSql()
	if err != nil {
		return rescode.Failed(err)
	}

	result, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return rescode.Failed(err)
	}

	if result.RowsAffected() == 0 {
		return rescode.SurveyNotFound(errors.New("survey not found"))
	}

	return nil
}

// DeleteSurvey deletes a survey
func (r *SurveyRepository) DeleteSurvey(ctx context.Context, id string) error {
	// Parse UUID from string
	surveyUUID, err := parseUUID(id)
	if err != nil {
		return rescode.IDInvalid(err)
	}

	query := r.sb.Delete("surveys").Where(squirrel.Eq{"id": surveyUUID})

	sql, args, err := query.ToSql()
	if err != nil {
		return rescode.Failed(err)
	}

	result, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return rescode.Failed(err)
	}

	if result.RowsAffected() == 0 {
		return rescode.SurveyNotFound(errors.New("survey not found"))
	}

	return nil
}

// ListSurveys lists all surveys
func (r *SurveyRepository) ListSurveys(ctx context.Context) ([]*entity.Survey, error) {
	query := r.sb.Select("id", "slug", "title", "description", "questions", "min_completion_time_min", "created_at", "updated_at", "created_by", "tags", "finishes_at").
		From("surveys").
		OrderBy("created_at DESC")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, rescode.Failed(err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, rescode.Failed(err)
	}
	defer rows.Close()

	var surveys []*entity.Survey

	for rows.Next() {
		var survey entity.Survey
		var questionsJSON []byte
		var tagsJSON []byte
		var dbUUID uuid.UUID

		err := rows.Scan(
			&dbUUID,
			&survey.Slug,
			&survey.Title,
			&survey.Description,
			&questionsJSON,
			&survey.MinCompletionTimeMin,
			&survey.CreatedAt,
			&survey.UpdatedAt,
			&survey.CreatedBy,
			&tagsJSON,
			&survey.FinishesAt,
		)
		if err != nil {
			return nil, rescode.Failed(err)
		}

		survey.ID = dbUUID.String()

		err = json.Unmarshal(questionsJSON, &survey.Questions)
		if err != nil {
			return nil, rescode.Failed(err)
		}

		err = json.Unmarshal(tagsJSON, &survey.Tags)
		if err != nil {
			return nil, rescode.Failed(err)
		}

		surveys = append(surveys, &survey)
	}

	if err := rows.Err(); err != nil {
		return nil, rescode.Failed(err)
	}

	return surveys, nil
}

// CreateSurveyResponse creates a new survey response
func (r *SurveyRepository) CreateSurveyResponse(ctx context.Context, response *entity.SurveyResponse) error {
	answersJSON, err := json.Marshal(response.Answers)
	if err != nil {
		return rescode.Failed(err)
	}

	// Parse UUIDs from strings
	responseUUID, err := parseUUID(response.ID)
	if err != nil {
		return rescode.IDInvalid(err)
	}

	var userUUID *uuid.UUID
	if response.UserID != "" {
		parsedUserUUID, err := parseUUID(response.UserID)
		if err != nil {
			return rescode.IDInvalid(err)
		}
		userUUID = &parsedUserUUID
	}

	survey, err := r.GetSurveyBySlug(ctx, response.SurveySlug)
	if err != nil {
		return rescode.SurveyNotFound(err)
	}

	query := r.sb.Insert("survey_responses").
		Columns("id", "survey_id", "user_id", "answers", "started_at", "completed_at", "ip_address", "user_agent", "is_completed", "is_anonymous", "created_at", "updated_at")

	if userUUID != nil {
		query = query.Values(
			responseUUID,
			survey.ID,
			userUUID,
			answersJSON,
			response.StartedAt,
			response.CompletedAt,
			response.IPAddress,
			response.UserAgent,
			response.IsCompleted,
			response.IsAnonymous,
			response.CreatedAt,
			response.UpdatedAt,
		)
	} else {
		query = query.Values(
			responseUUID,
			survey.ID,
			nil,
			answersJSON,
			response.StartedAt,
			response.CompletedAt,
			response.IPAddress,
			response.UserAgent,
			response.IsCompleted,
			response.IsAnonymous,
			response.CreatedAt,
			response.UpdatedAt,
		)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return rescode.Failed(err)
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return rescode.Failed(err)
	}

	return nil
}

// GetSurveyResponseByID gets a survey response by ID
func (r *SurveyRepository) GetSurveyResponseByID(ctx context.Context, id string) (*entity.SurveyResponse, error) {
	// Parse UUID from string
	responseUUID, err := parseUUID(id)
	if err != nil {
		return nil, rescode.IDInvalid(err)
	}

	query := r.sb.Select("id", "survey_id", "user_id", "answers", "started_at", "completed_at", "ip_address", "user_agent", "is_completed", "is_anonymous", "created_at", "updated_at").
		From("survey_responses").
		Where(squirrel.Eq{"id": responseUUID})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, rescode.Failed(err)
	}

	var response entity.SurveyResponse
	var answersJSON []byte
	var completedAt *time.Time
	var responseIDUUID uuid.UUID
	var surveySlug string
	var userID *uuid.UUID

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&responseIDUUID,
		&surveySlug,
		&userID,
		&answersJSON,
		&response.StartedAt,
		&completedAt,
		&response.IPAddress,
		&response.UserAgent,
		&response.IsCompleted,
		&response.IsAnonymous,
		&response.CreatedAt,
		&response.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, rescode.SurveyNotFound(err)
		}
		return nil, rescode.Failed(err)
	}

	response.ID = responseIDUUID.String()
	response.SurveySlug = surveySlug
	if userID != nil {
		response.UserID = userID.String()
	}
	response.CompletedAt = completedAt

	err = json.Unmarshal(answersJSON, &response.Answers)
	if err != nil {
		return nil, rescode.Failed(err)
	}

	return &response, nil
}

// UpdateSurveyResponse updates a survey response
func (r *SurveyRepository) UpdateSurveyResponse(ctx context.Context, response *entity.SurveyResponse) error {
	answersJSON, err := json.Marshal(response.Answers)
	if err != nil {
		return rescode.Failed(err)
	}

	// Parse UUID from string
	responseUUID, err := parseUUID(response.ID)
	if err != nil {
		return rescode.IDInvalid(err)
	}

	query := r.sb.Update("survey_responses").
		Set("answers", answersJSON).
		Set("completed_at", response.CompletedAt).
		Set("is_completed", response.IsCompleted).
		Set("updated_at", response.UpdatedAt).
		Where(squirrel.Eq{"id": responseUUID})

	sql, args, err := query.ToSql()
	if err != nil {
		return rescode.Failed(err)
	}

	result, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return rescode.Failed(err)
	}

	if result.RowsAffected() == 0 {
		return rescode.SurveyNotFound(errors.New("survey response not found"))
	}

	return nil
}

// ListSurveyResponses lists all responses for a survey
func (r *SurveyRepository) ListSurveyResponses(ctx context.Context, surveySlug string) ([]*entity.SurveyResponse, error) {
	query := r.sb.Select("id", "survey_slug", "user_id", "answers", "started_at", "completed_at", "ip_address", "user_agent", "is_completed", "is_anonymous", "created_at", "updated_at").
		From("survey_responses").
		Where(squirrel.Eq{"survey_slug": surveySlug}).
		OrderBy("created_at DESC")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, rescode.Failed(err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, rescode.Failed(err)
	}
	defer rows.Close()

	var responses []*entity.SurveyResponse

	for rows.Next() {
		var response entity.SurveyResponse
		var answersJSON []byte
		var completedAt *time.Time
		var responseIDUUID uuid.UUID
		var surveySlug string
		var userIDUUID *uuid.UUID

		err := rows.Scan(
			&responseIDUUID,
			&surveySlug,
			&userIDUUID,
			&answersJSON,
			&response.StartedAt,
			&completedAt,
			&response.IPAddress,
			&response.UserAgent,
			&response.IsCompleted,
			&response.IsAnonymous,
			&response.CreatedAt,
			&response.UpdatedAt,
		)
		if err != nil {
			return nil, rescode.Failed(err)
		}

		response.ID = responseIDUUID.String()
		response.SurveySlug = surveySlug
		if userIDUUID != nil {
			response.UserID = userIDUUID.String()
		}
		response.CompletedAt = completedAt

		err = json.Unmarshal(answersJSON, &response.Answers)
		if err != nil {
			return nil, rescode.Failed(err)
		}

		responses = append(responses, &response)
	}

	if err := rows.Err(); err != nil {
		return nil, rescode.Failed(err)
	}

	return responses, nil
}

// GetSurveyResponseStats gets statistics for a survey
func (r *SurveyRepository) GetSurveyResponseStats(ctx context.Context, surveyID string) (map[string]interface{}, error) {
	// Get the survey to access its questions
	survey, err := r.GetSurveyByID(ctx, surveyID)
	if err != nil {
		return nil, rescode.Failed(err)
	}

	// Parse UUID from string
	surveyUUID, err := parseUUID(surveyID)
	if err != nil {
		return nil, rescode.IDInvalid(err)
	}

	// Get all completed responses for the survey
	query := r.sb.Select("answers").
		From("survey_responses").
		Where(squirrel.Eq{"survey_id": surveyUUID, "is_completed": true})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, rescode.Failed(err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, rescode.Failed(err)
	}
	defer rows.Close()

	// Initialize stats
	stats := map[string]interface{}{
		"totalResponses": 0,
		"questions":      make(map[string]interface{}),
	}

	// Initialize question stats
	for _, question := range survey.Questions {
		switch question.Type {
		case entity.QuestionTypeNumber:
			stats["questions"].(map[string]interface{})[question.ID] = map[string]interface{}{
				"type":   question.Type,
				"text":   question.Text,
				"count":  0,
				"sum":    0.0,
				"avg":    0.0,
				"min":    nil,
				"max":    nil,
				"median": 0.0,
			}
		case entity.QuestionTypeSelect:
			options := make(map[string]int)
			for _, option := range question.Options {
				options[option.ID] = 0
			}
			stats["questions"].(map[string]interface{})[question.ID] = map[string]interface{}{
				"type":    question.Type,
				"text":    question.Text,
				"count":   0,
				"options": options,
			}
		case entity.QuestionTypeMultiSelect:
			options := make(map[string]int)
			for _, option := range question.Options {
				options[option.ID] = 0
			}
			stats["questions"].(map[string]interface{})[question.ID] = map[string]interface{}{
				"type":    question.Type,
				"text":    question.Text,
				"count":   0,
				"options": options,
			}
		case entity.QuestionTypeText:
			stats["questions"].(map[string]interface{})[question.ID] = map[string]interface{}{
				"type":  question.Type,
				"text":  question.Text,
				"count": 0,
			}
		case entity.QuestionTypeBoolean:
			stats["questions"].(map[string]interface{})[question.ID] = map[string]interface{}{
				"type":  question.Type,
				"text":  question.Text,
				"count": 0,
				"true":  0,
				"false": 0,
			}
		}
	}

	// Process responses
	responseCount := 0
	for rows.Next() {
		var answersJSON []byte
		err := rows.Scan(&answersJSON)
		if err != nil {
			return nil, rescode.Failed(err)
		}

		var answers map[string]interface{}
		err = json.Unmarshal(answersJSON, &answers)
		if err != nil {
			return nil, rescode.Failed(err)
		}

		responseCount++

		// Process each answer
		for questionID, answer := range answers {
			questionStats, ok := stats["questions"].(map[string]interface{})[questionID]
			if !ok {
				continue // Skip if question not found in survey
			}

			// Update question stats based on question type
			switch questionStats.(map[string]interface{})["type"] {
			case entity.QuestionTypeNumber:
				if answer == nil {
					continue
				}

				// Convert answer to float64
				var value float64
				switch v := answer.(type) {
				case float64:
					value = v
				case int:
					value = float64(v)
				case string:
					// Try to parse string as float
					var parseErr error
					value, parseErr = parseFloat(v)
					if parseErr != nil {
						continue
					}
				default:
					continue
				}

				qs := questionStats.(map[string]interface{})
				qs["count"] = qs["count"].(int) + 1
				qs["sum"] = qs["sum"].(float64) + value

				// Update min/max
				if qs["min"] == nil || value < qs["min"].(float64) {
					qs["min"] = value
				}
				if qs["max"] == nil || value > qs["max"].(float64) {
					qs["max"] = value
				}

				// Calculate average
				if qs["count"].(int) > 0 {
					qs["avg"] = qs["sum"].(float64) / float64(qs["count"].(int))
				}

			case entity.QuestionTypeSelect:
				if answer == nil {
					continue
				}

				qs := questionStats.(map[string]interface{})
				qs["count"] = qs["count"].(int) + 1

				// Convert answer to string
				var optionID string
				switch v := answer.(type) {
				case string:
					optionID = v
				default:
					continue
				}

				// Update option count
				options := qs["options"].(map[string]int)
				if _, ok := options[optionID]; ok {
					options[optionID]++
				}

			case entity.QuestionTypeMultiSelect:
				if answer == nil {
					continue
				}

				qs := questionStats.(map[string]interface{})
				qs["count"] = qs["count"].(int) + 1

				// Convert answer to []string
				var optionIDs []string
				switch v := answer.(type) {
				case []interface{}:
					for _, item := range v {
						if str, ok := item.(string); ok {
							optionIDs = append(optionIDs, str)
						}
					}
				default:
					continue
				}

				// Update option counts
				options := qs["options"].(map[string]int)
				for _, optionID := range optionIDs {
					if _, ok := options[optionID]; ok {
						options[optionID]++
					}
				}

			case entity.QuestionTypeText:
				if answer == nil {
					continue
				}

				qs := questionStats.(map[string]interface{})
				qs["count"] = qs["count"].(int) + 1

			case entity.QuestionTypeBoolean:
				if answer == nil {
					continue
				}

				qs := questionStats.(map[string]interface{})
				qs["count"] = qs["count"].(int) + 1

				// Convert answer to bool
				var value bool
				switch v := answer.(type) {
				case bool:
					value = v
				case string:
					value = v == "true"
				default:
					continue
				}

				if value {
					qs["true"] = qs["true"].(int) + 1
				} else {
					qs["false"] = qs["false"].(int) + 1
				}
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, rescode.Failed(err)
	}

	stats["totalResponses"] = responseCount

	return stats, nil
}

// Helper function to parse a string as a float
func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

// GetSurveyBySlug gets a survey by slug
func (r *SurveyRepository) GetSurveyBySlug(ctx context.Context, slug string) (*entity.Survey, error) {
	query := r.sb.Select("id", "slug", "title", "description", "questions", "min_completion_time_min", "created_at", "updated_at", "created_by", "tags", "finishes_at").
		From("surveys").
		Where(squirrel.Eq{"slug": slug})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, rescode.Failed(err)
	}

	var survey entity.Survey
	var questionsJSON []byte
	var tagsJSON []byte
	var dbUUID uuid.UUID

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&dbUUID,
		&survey.Slug,
		&survey.Title,
		&survey.Description,
		&questionsJSON,
		&survey.MinCompletionTimeMin,
		&survey.CreatedAt,
		&survey.UpdatedAt,
		&survey.CreatedBy,
		&tagsJSON,
		&survey.FinishesAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, rescode.SurveyNotFound(err)
		}
		return nil, rescode.Failed(err)
	}

	survey.ID = dbUUID.String()

	err = json.Unmarshal(questionsJSON, &survey.Questions)
	if err != nil {
		return nil, rescode.Failed(err)
	}

	err = json.Unmarshal(tagsJSON, &survey.Tags)
	if err != nil {
		return nil, rescode.Failed(err)
	}

	return &survey, nil
}

// CountSurveyParticipants counts the number of completed responses for a survey
func (r *SurveyRepository) CountSurveyParticipants(ctx context.Context, surveyID string) (int, error) {
	// Parse UUID from string
	surveyUUID, err := parseUUID(surveyID)
	if err != nil {
		return 0, rescode.IDInvalid(err)
	}

	query := r.sb.Select("COUNT(*)").
		From("survey_responses").
		Where(squirrel.Eq{"survey_id": surveyUUID, "is_completed": true})

	sql, args, err := query.ToSql()
	if err != nil {
		return 0, rescode.Failed(err)
	}

	var count int
	err = r.db.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		return 0, rescode.Failed(err)
	}

	return count, nil
}

// ListSurveysWithPagination lists surveys with pagination and filtering
func (r *SurveyRepository) ListSurveysWithPagination(ctx context.Context, req *list.PagiRequest, filterReq *entity.SurveyListRequest) (*list.PagiResponse[*entity.SurveyListItem], error) {
	// Set default values for pagination
	req.Default()

	// Build the base query
	baseQuery := r.sb.Select("surveys.id", "surveys.slug", "surveys.title", "surveys.description", "surveys.min_completion_time_min", "surveys.created_at", "surveys.updated_at", "surveys.created_by", "surveys.tags", "surveys.finishes_at", "COUNT(survey_responses.id) AS participants").
		LeftJoin("survey_responses ON surveys.id = survey_responses.survey_id").
		From("surveys").
		GroupBy("surveys.id")

	// Apply tag filtering if provided
	if filterReq != nil && len(filterReq.Tag) > 0 {
		baseQuery = baseQuery.Where("tags @> ?", `["`+filterReq.Tag+`"]`)
	}

	// Filter out expired surveys if requested
	if filterReq != nil && filterReq.HideExpired {
		baseQuery = baseQuery.Where("finishes_at IS NULL OR finishes_at > NOW()")
	}

	// Apply text search if provided
	if filterReq != nil && filterReq.Search != "" {
		// Create a tsquery from the search term
		tsQuery := filterReq.Search
		// Replace spaces with & for AND operations in the search
		tsQuery = strings.ReplaceAll(tsQuery, " ", " & ")
		// Add :* to each word for prefix matching
		tsQuery = strings.ReplaceAll(tsQuery, " & ", ":* & ")
		tsQuery = tsQuery + ":*"

		// Apply the text search condition using to_tsvector and to_tsquery
		baseQuery = baseQuery.Where("to_tsvector('simple', surveys.title || ' ' || surveys.description) @@ to_tsquery('simple', ?)", tsQuery)
	}

	// Apply sorting
	if filterReq != nil && filterReq.Sort != "" {
		switch filterReq.Sort {
		case "created_at_asc":
			baseQuery = baseQuery.OrderBy("surveys.created_at ASC")
		case "created_at_desc":
			baseQuery = baseQuery.OrderBy("surveys.created_at DESC")
		case "most_participants":
			baseQuery = baseQuery.OrderBy("participants DESC")
		case "finishes_at_asc":
			baseQuery = baseQuery.OrderBy("surveys.finishes_at ASC NULLS LAST")
		case "finishes_at_desc":
			baseQuery = baseQuery.OrderBy("surveys.finishes_at DESC NULLS LAST")
		default:
			baseQuery = baseQuery.OrderBy("surveys.created_at DESC") // Default sorting
		}
	} else {
		baseQuery = baseQuery.OrderBy("surveys.created_at DESC") // Default sorting
	}

	baseQuery = baseQuery.Limit(*req.Limit).Offset(req.Offset())

	sql, args, err := baseQuery.ToSql()
	if err != nil {
		return nil, rescode.Failed(err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, rescode.Failed(err)
	}
	defer rows.Close()

	var surveyItems []*entity.SurveyListItem

	for rows.Next() {
		var survey entity.Survey
		var tagsJSON []byte
		var participants int
		var dbUUID uuid.UUID

		err := rows.Scan(
			&dbUUID,
			&survey.Slug,
			&survey.Title,
			&survey.Description,
			&survey.MinCompletionTimeMin,
			&survey.CreatedAt,
			&survey.UpdatedAt,
			&survey.CreatedBy,
			&tagsJSON,
			&survey.FinishesAt,
			&participants,
		)
		if err != nil {
			return nil, rescode.Failed(err)
		}

		survey.ID = dbUUID.String()

		err = json.Unmarshal(tagsJSON, &survey.Tags)
		if err != nil {
			return nil, rescode.Failed(err)
		}

		// Convert to list item
		listItem := survey.ToListItem(participants)
		listItem.FinishesAt = survey.FinishesAt

		surveyItems = append(surveyItems, listItem)
	}

	if err := rows.Err(); err != nil {
		return nil, rescode.Failed(err)
	}

	// Create pagination response
	response := &list.PagiResponse[*entity.SurveyListItem]{
		Page:  *req.Page,
		Limit: *req.Limit,
		List:  surveyItems,
	}

	return response, nil
}
