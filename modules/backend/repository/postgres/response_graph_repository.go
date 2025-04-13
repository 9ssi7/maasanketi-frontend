package postgres

import (
	"context"
	"encoding/json"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mstrYoda/maasanketi.co/domain/resgraph"
	"github.com/mstrYoda/maasanketi.co/pkg/rescode"
)

type ResponseGraphRepository struct {
	db *pgxpool.Pool
	sb squirrel.StatementBuilderType
}

func NewResponseGraphRepository(db *pgxpool.Pool) *ResponseGraphRepository {
	return &ResponseGraphRepository{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *ResponseGraphRepository) Create(ctx context.Context, responseGraph *resgraph.ResponseGraph) error {
	query := r.sb.Insert("response_graphs").
		Columns("id", "survey_id", "graph_id", "content", "created_at", "updated_at")

	contentJSON, err := json.Marshal(responseGraph.Content)
	if err != nil {
		return rescode.Failed(err)
	}

	query = query.Values(
		responseGraph.ID,
		responseGraph.SurveyID,
		responseGraph.GraphID,
		contentJSON,
		responseGraph.CreatedAt,
		responseGraph.UpdatedAt,
	)

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

func (r *ResponseGraphRepository) FindByID(ctx context.Context, id string) (*resgraph.ResponseGraph, error) {
	query := r.sb.Select("id", "survey_id", "graph_id", "content", "created_at", "updated_at").
		From("response_graphs").
		Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, rescode.Failed(err)
	}

	var responseGraph resgraph.ResponseGraph
	var contentJSON []byte

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&responseGraph.ID,
		&responseGraph.SurveyID,
		&responseGraph.GraphID,
		&contentJSON,
		&responseGraph.CreatedAt,
		&responseGraph.UpdatedAt,
	)
	if err != nil {
		return nil, rescode.Failed(err)
	}

	err = json.Unmarshal(contentJSON, &responseGraph.Content)
	if err != nil {
		return nil, rescode.Failed(err)
	}

	return &responseGraph, nil
}

func (r *ResponseGraphRepository) Delete(ctx context.Context, id string) error {
	query := r.sb.Delete("response_graphs").
		Where(squirrel.Eq{"id": id})

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

func (r *ResponseGraphRepository) ListBySurveyID(ctx context.Context, surveyID string) ([]*resgraph.ViewList, error) {
	query := r.sb.Select("response_graphs.id", "response_graphs.survey_id", "response_graphs.graph_id", "response_graphs.content", "response_graphs.created_at", "response_graphs.updated_at", "graphs.title", "graphs.description", "graphs.kind", "graphs.content").
		From("response_graphs").
		LeftJoin("graphs ON response_graphs.graph_id = graphs.id").
		Where("response_graphs.survey_id = ?", surveyID).
		OrderBy("response_graphs.created_at DESC")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, rescode.Failed(err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, rescode.Failed(err)
	}
	defer rows.Close()

	var responseGraphs []*resgraph.ViewList

	for rows.Next() {
		var responseGraph resgraph.ViewList
		var contentJSON []byte

		err = rows.Scan(
			&responseGraph.ID,
			&responseGraph.SurveyID,
			&responseGraph.GraphID,
			&contentJSON,
			&responseGraph.CreatedAt,
			&responseGraph.UpdatedAt,
			&responseGraph.Title,
			&responseGraph.Description,
			&responseGraph.Kind,
			&responseGraph.GraphContent,
		)
		if err != nil {
			return nil, rescode.Failed(err)
		}

		err = json.Unmarshal(contentJSON, &responseGraph.Content)
		if err != nil {
			return nil, rescode.Failed(err)
		}

		responseGraphs = append(responseGraphs, &responseGraph)
	}

	if err := rows.Err(); err != nil {
		return nil, rescode.Failed(err)
	}

	return responseGraphs, nil
}

func (r *ResponseGraphRepository) ListByGraphID(ctx context.Context, graphID string) ([]*resgraph.ResponseGraph, error) {
	query := r.sb.Select("id", "survey_id", "graph_id", "content", "created_at", "updated_at").
		From("response_graphs").
		Where(squirrel.Eq{"graph_id": graphID})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, rescode.Failed(err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, rescode.Failed(err)
	}
	defer rows.Close()

	var responseGraphs []*resgraph.ResponseGraph

	for rows.Next() {
		var responseGraph resgraph.ResponseGraph
		var contentJSON []byte

		err = rows.Scan(
			&responseGraph.ID,
			&responseGraph.SurveyID,
			&responseGraph.GraphID,
			&contentJSON,
			&responseGraph.CreatedAt,
			&responseGraph.UpdatedAt,
		)
		if err != nil {
			return nil, rescode.Failed(err)
		}

		err = json.Unmarshal(contentJSON, &responseGraph.Content)
		if err != nil {
			return nil, rescode.Failed(err)
		}

		responseGraphs = append(responseGraphs, &responseGraph)
	}

	if err := rows.Err(); err != nil {
		return nil, rescode.Failed(err)
	}

	return responseGraphs, nil
}

func (r *ResponseGraphRepository) Update(ctx context.Context, responseGraph *resgraph.ResponseGraph) error {
	query := r.sb.Update("response_graphs").
		Set("content", responseGraph.Content).
		Set("updated_at", responseGraph.UpdatedAt).
		Where("id = ?", responseGraph.ID)

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
