package postgres

import (
	"context"
	"encoding/json"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mstrYoda/maasanketi.co/domain/graph"
	"github.com/mstrYoda/maasanketi.co/pkg/rescode"
)

type GraphRepository struct {
	db *pgxpool.Pool
	sb squirrel.StatementBuilderType
}

func NewGraphRepository(db *pgxpool.Pool) *GraphRepository {
	return &GraphRepository{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *GraphRepository) Create(ctx context.Context, graph *graph.Graph) error {
	query := r.sb.Insert("graphs").
		Columns("survey_id", "title", "description", "kind", "content", "created_at", "updated_at")

	contentJSON, err := json.Marshal(graph.Content)
	if err != nil {
		return rescode.Failed(err)
	}

	query = query.Values(
		uuid.MustParse(graph.SurveyID),
		graph.Title,
		graph.Description,
		graph.Kind,
		contentJSON,
		graph.CreatedAt,
		graph.UpdatedAt,
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

func (r *GraphRepository) FindByID(ctx context.Context, id string) (*graph.Graph, error) {
	query := r.sb.Select("id", "survey_id", "title", "description", "kind", "content", "created_at", "updated_at").
		From("graphs").
		Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, rescode.Failed(err)
	}

	var graph graph.Graph
	var contentJSON []byte

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&graph.ID,
		&graph.SurveyID,
		&graph.Title,
		&graph.Description,
		&graph.Kind,
		&contentJSON,
		&graph.CreatedAt,
		&graph.UpdatedAt,
	)
	if err != nil {
		return nil, rescode.Failed(err)
	}

	err = json.Unmarshal(contentJSON, &graph.Content)
	if err != nil {
		return nil, rescode.Failed(err)
	}

	return &graph, nil
}

func (r *GraphRepository) Delete(ctx context.Context, id string) error {
	query := r.sb.Delete("graphs").
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

func (r *GraphRepository) ListBySurveyID(ctx context.Context, surveyID string) ([]*graph.Graph, error) {
	query := r.sb.Select("id", "survey_id", "title", "description", "kind", "content", "created_at", "updated_at").
		From("graphs").
		Where(squirrel.Eq{"survey_id": surveyID})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, rescode.Failed(err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, rescode.Failed(err)
	}
	defer rows.Close()

	var graphs []*graph.Graph

	for rows.Next() {
		var graph graph.Graph
		var contentJSON []byte

		err = rows.Scan(
			&graph.ID,
			&graph.SurveyID,
			&graph.Title,
			&graph.Description,
			&graph.Kind,
			&contentJSON,
			&graph.CreatedAt,
			&graph.UpdatedAt,
		)
		if err != nil {
			return nil, rescode.Failed(err)
		}

		err = json.Unmarshal(contentJSON, &graph.Content)
		if err != nil {
			return nil, rescode.Failed(err)
		}

		graphs = append(graphs, &graph)
	}

	if err := rows.Err(); err != nil {
		return nil, rescode.Failed(err)
	}

	return graphs, nil
}
