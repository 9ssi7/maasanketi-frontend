package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mstrYoda/maasanketi.co/domain/response"
	"github.com/mstrYoda/maasanketi.co/pkg/rescode"
)

type ResponseRepository struct {
	db *pgxpool.Pool
	sb squirrel.StatementBuilderType
}

func NewResponseRepository(db *pgxpool.Pool) *ResponseRepository {
	return &ResponseRepository{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// ResponseCreate creates a new survey response
func (r *ResponseRepository) Create(ctx context.Context, response *response.Response) error {
	answersJSON, err := json.Marshal(response.Answers)
	if err != nil {
		return rescode.Failed(err)
	}

	query := r.sb.Insert("survey_responses").
		Columns("id", "survey_id", "user_id", "answers", "started_at", "completed_at", "ip_address", "user_agent", "is_completed", "is_anonymous", "created_at", "updated_at")

	if response.UserID != "" {
		query = query.Values(
			uuid.MustParse(response.ID),
			response.SurveyID,
			uuid.MustParse(response.UserID),
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
			uuid.MustParse(response.ID),
			response.SurveyID,
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

// ResponseFindByID gets a survey response by ID
func (r *ResponseRepository) FindByID(ctx context.Context, id string) (*response.Response, error) {
	query := r.sb.Select("id", "survey_id", "user_id", "answers", "started_at", "completed_at", "ip_address", "user_agent", "is_completed", "is_anonymous", "created_at", "updated_at").
		From("survey_responses").
		Where(squirrel.Eq{"id": uuid.MustParse(id)})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, rescode.Failed(err)
	}

	var response response.Response
	var answersJSON []byte
	var completedAt *time.Time
	var responseID uuid.UUID
	var surveyID string
	var userID *uuid.UUID

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&responseID,
		&surveyID,
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

	response.ID = responseID.String()
	response.SurveyID = surveyID
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

// ResponseUpdate updates a survey response
func (r *ResponseRepository) Update(ctx context.Context, response *response.Response) error {
	answersJSON, err := json.Marshal(response.Answers)
	if err != nil {
		return rescode.Failed(err)
	}

	query := r.sb.Update("survey_responses").
		Set("answers", answersJSON).
		Set("completed_at", response.CompletedAt).
		Set("is_completed", response.IsCompleted).
		Set("updated_at", response.UpdatedAt).
		Where(squirrel.Eq{"id": uuid.MustParse(response.ID)})

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

// ResponseList lists all responses for a survey
func (r *ResponseRepository) List(ctx context.Context, surveyID string) ([]*response.Response, error) {
	query := r.sb.Select("id", "survey_id", "user_id", "answers", "started_at", "completed_at", "ip_address", "user_agent", "is_completed", "is_anonymous", "created_at", "updated_at").
		From("survey_responses").
		Where(squirrel.Eq{"survey_id": uuid.MustParse(surveyID)}).
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

	var responses []*response.Response

	for rows.Next() {
		var response response.Response
		var answersJSON []byte
		var completedAt *time.Time
		var responseIDUUID uuid.UUID
		var surveyID string
		var userIDUUID *uuid.UUID

		err := rows.Scan(
			&responseIDUUID,
			&surveyID,
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
		response.SurveyID = surveyID
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
