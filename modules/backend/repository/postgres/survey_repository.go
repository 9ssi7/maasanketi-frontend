package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mstrYoda/maasanketi.co/domain/survey"
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

// CreateSuCreatervey creates a new survey
func (r *SurveyRepository) Create(ctx context.Context, survey *survey.Survey) error {
	questionsJSON, err := json.Marshal(survey.Questions)
	if err != nil {
		return rescode.Failed(err)
	}

	tagsJSON, err := json.Marshal(survey.Tags)
	if err != nil {
		return rescode.Failed(err)
	}

	// Generate slug if not provided
	if survey.Slug == "" {
		survey.GenerateSlug()
	}

	query := r.sb.Insert("surveys").
		Columns("id", "slug", "title", "description", "questions", "min_completion_time_min", "created_at", "updated_at", "created_by", "tags", "finishes_at").
		Values(uuid.MustParse(survey.ID), survey.Slug, survey.Title, survey.Description, questionsJSON, survey.MinCompletionTimeMin, survey.CreatedAt, survey.UpdatedAt, survey.CreatedBy, tagsJSON, survey.FinishesAt)

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

// FindByID gets a survey by ID
func (r *SurveyRepository) FindByID(ctx context.Context, id string) (*survey.Survey, error) {
	query := r.sb.Select("id", "slug", "title", "description", "questions", "min_completion_time_min", "created_at", "updated_at", "created_by", "tags", "finishes_at").
		From("surveys").
		Where(squirrel.Eq{"id": uuid.MustParse(id)})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, rescode.Failed(err)
	}

	var survey survey.Survey
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

// Update updates a survey
func (r *SurveyRepository) Update(ctx context.Context, survey *survey.Survey) error {
	questionsJSON, err := json.Marshal(survey.Questions)
	if err != nil {
		return rescode.Failed(err)
	}

	tagsJSON, err := json.Marshal(survey.Tags)
	if err != nil {
		return rescode.Failed(err)
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
		Where(squirrel.Eq{"id": uuid.MustParse(survey.ID)})

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

// Delete deletes a survey
func (r *SurveyRepository) Delete(ctx context.Context, id string) error {
	query := r.sb.Delete("surveys").Where(squirrel.Eq{"id": uuid.MustParse(id)})

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

// Helper function to parse a string as a float
func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

// FindBySlug gets a survey by slug
func (r *SurveyRepository) FindBySlug(ctx context.Context, slug string) (*survey.Survey, error) {
	query := r.sb.Select("id", "slug", "title", "description", "questions", "min_completion_time_min", "created_at", "updated_at", "created_by", "tags", "finishes_at").
		From("surveys").
		Where(squirrel.Eq{"slug": slug})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, rescode.Failed(err)
	}

	var survey survey.Survey
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

// CountParticipants counts the number of completed responses for a survey
func (r *SurveyRepository) CountParticipants(ctx context.Context, surveyID string) (int, error) {
	query := r.sb.Select("COUNT(*)").
		From("survey_responses").
		Where(squirrel.Eq{"survey_id": uuid.MustParse(surveyID), "is_completed": true})

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

// FindCompletedSurveys finds completed surveys
func (r *SurveyRepository) FindCompletedSurveys(ctx context.Context, page, limit uint64) ([]*survey.Survey, error) {
	query := r.sb.Select("surveys.id", "surveys.slug", "surveys.title", "surveys.description", "surveys.min_completion_time_min", "surveys.created_at", "surveys.updated_at", "surveys.created_by", "surveys.tags", "surveys.finishes_at", "COUNT(survey_responses.id) AS participants").
		From("surveys").
		LeftJoin("survey_responses ON surveys.id = survey_responses.survey_id").
		LeftJoin("response_graphs ON surveys.id = response_graphs.survey_id").
		Where("surveys.finishes_at < NOW()").
		Where("response_graphs.id IS NULL").
		OrderBy("surveys.created_at DESC").
		GroupBy("surveys.id").
		Limit(limit).Offset(page * limit)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, rescode.Failed(err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, rescode.Failed(err)
	}
	defer rows.Close()

	var surveys []*survey.Survey

	for rows.Next() {
		var survey survey.Survey
		var participants int
		var dbUUID uuid.UUID
		var tagsJSON []byte
		err := rows.Scan(&dbUUID, &survey.Slug, &survey.Title, &survey.Description, &survey.MinCompletionTimeMin, &survey.CreatedAt, &survey.UpdatedAt, &survey.CreatedBy, &tagsJSON, &survey.FinishesAt, &participants)
		if err != nil {
			return nil, rescode.Failed(err)
		}

		survey.ID = dbUUID.String()

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

// ViewDetail gets a survey by slug
func (r *SurveyRepository) ViewDetail(ctx context.Context, slug string) (*survey.ViewDetail, error) {
	query := r.sb.Select("surveys.id", "surveys.slug", "surveys.title", "surveys.description", "surveys.questions", "surveys.min_completion_time_min", "surveys.created_at", "surveys.updated_at", "surveys.created_by", "surveys.tags", "surveys.finishes_at", "COUNT(survey_responses.id) AS participants").
		From("surveys").
		LeftJoin("survey_responses ON surveys.id = survey_responses.survey_id").
		Where(squirrel.Eq{"slug": slug})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, rescode.Failed(err)
	}

	var s survey.Survey
	var questionsJSON []byte
	var tagsJSON []byte
	var dbUUID uuid.UUID
	var participants int

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&dbUUID,
		&s.Slug,
		&s.Title,
		&s.Description,
		&questionsJSON,
		&s.MinCompletionTimeMin,
		&s.CreatedAt,
		&s.UpdatedAt,
		&s.CreatedBy,
		&tagsJSON,
		&s.FinishesAt,
		&participants,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, rescode.SurveyNotFound(err)
		}
		return nil, rescode.Failed(err)
	}

	s.ID = dbUUID.String()

	err = json.Unmarshal(questionsJSON, &s.Questions)
	if err != nil {
		return nil, rescode.Failed(err)
	}

	err = json.Unmarshal(tagsJSON, &s.Tags)
	if err != nil {
		return nil, rescode.Failed(err)
	}

	viewDetail := &survey.ViewDetail{
		ViewList:  *s.ToListView(participants),
		Questions: s.Questions,
	}
	return viewDetail, nil
}

// List lists surveys with pagination and filtering
func (r *SurveyRepository) List(ctx context.Context, req *list.PagiRequest, filterReq *survey.SurveyListFilters) (*list.PagiResponse[*survey.ViewList], error) {
	// Set default values for pagination
	req.Default()

	// Build the base query
	baseQuery := r.sb.Select("surveys.id", "surveys.slug", "surveys.title", "surveys.description", "surveys.min_completion_time_min", "surveys.created_at", "surveys.updated_at", "surveys.created_by", "surveys.tags", "surveys.finishes_at", "COUNT(survey_responses.id) AS participants").
		From("surveys").
		LeftJoin("survey_responses ON surveys.id = survey_responses.survey_id").
		GroupBy("surveys.id")

	// Apply tag filtering if provided
	if filterReq != nil && len(filterReq.Tag) > 0 {
		baseQuery = baseQuery.Where("tags @> ?", `["`+filterReq.Tag+`"]`)
	}

	// Filter out expired surveys if requested
	if filterReq != nil && filterReq.HideExpired != nil {
		if *filterReq.HideExpired == "true" {
			baseQuery = baseQuery.Where("surveys.finishes_at IS NULL OR surveys.finishes_at > NOW()")
		} else {
			baseQuery = baseQuery.Where("surveys.finishes_at < NOW()")
		}
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

	var surveyItems []*survey.ViewList

	for rows.Next() {
		var survey survey.Survey
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
		listItem := survey.ToListView(participants)
		listItem.FinishesAt = survey.FinishesAt

		surveyItems = append(surveyItems, listItem)
	}

	if err := rows.Err(); err != nil {
		return nil, rescode.Failed(err)
	}

	// Create pagination response
	response := &list.PagiResponse[*survey.ViewList]{
		Page:  *req.Page,
		Limit: *req.Limit,
		List:  surveyItems,
	}

	return response, nil
}
