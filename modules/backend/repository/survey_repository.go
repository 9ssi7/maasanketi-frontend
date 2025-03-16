package repository

import (
	"context"

	"github.com/mstrYoda/maasanketi.co/entity"
)

// SurveyRepository defines the interface for survey data access
type SurveyRepository interface {
	// Survey operations
	CreateSurvey(ctx context.Context, survey *entity.Survey) error
	GetSurveyByID(ctx context.Context, id string) (*entity.Survey, error)
	GetSurveyBySlug(ctx context.Context, slug string) (*entity.Survey, error)
	UpdateSurvey(ctx context.Context, survey *entity.Survey) error
	DeleteSurvey(ctx context.Context, id string) error
	ListSurveys(ctx context.Context) ([]*entity.Survey, error)

	// SurveyResponse operations
	CreateSurveyResponse(ctx context.Context, response *entity.SurveyResponse) error
	GetSurveyResponseByID(ctx context.Context, id string) (*entity.SurveyResponse, error)
	UpdateSurveyResponse(ctx context.Context, response *entity.SurveyResponse) error
	ListSurveyResponses(ctx context.Context, surveyID string) ([]*entity.SurveyResponse, error)
	GetSurveyResponseStats(ctx context.Context, surveyID string) (map[string]interface{}, error)
}
