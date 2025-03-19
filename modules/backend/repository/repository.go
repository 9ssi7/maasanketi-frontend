package repository

import (
	"context"
	"errors"
	"time"

	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/mstrYoda/maasanketi.co/pkg/logger"
	"github.com/mstrYoda/maasanketi.co/repository/postgres"
)

type Repository struct {
	pool *pgxpool.Pool

	Survey   *postgres.SurveyRepository
	Response *postgres.ResponseRepository
}

func New() (*Repository, error) {
	dbUrl := os.Getenv("DB_CONN_STR")
	if dbUrl == "" {
		return nil, errors.New("DB_CONN_STR is not set")
	}
	config, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		return nil, errors.New("unable to parse config: " + err.Error())
	}

	config.MinConns = 5
	config.MaxConns = 10
	config.ConnConfig.Tracer = &queryTracer{
		log: log.Logger().Sugar(),
	}

	pool, err := pgxpool.NewWithConfig(
		context.Background(),
		config,
	)
	if err != nil {
		return nil, errors.New("unable to connect to database: " + err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		return nil, errors.New("unable to ping database: " + err.Error())
	}

	repo := &Repository{
		pool: pool,
	}

	// Initialize repositories
	repo.Survey = postgres.NewSurveyRepository(pool)
	repo.Response = postgres.NewResponseRepository(pool)
	return repo, nil
}

func (r *Repository) Close() {
	r.pool.Close()
}
