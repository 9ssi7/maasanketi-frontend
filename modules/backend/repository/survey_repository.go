package repository

import (
	"context"

	"github.com/mstrYoda/maasanketi.co/entity"
	"github.com/mstrYoda/maasanketi.co/pkg/list"
)

// SurveyRepository defines the interface for survey data access
type SurveyRepository interface {
	// Survey operations
	Create(ctx context.Context, survey *entity.Survey) error
	FindByID(ctx context.Context, id string) (*entity.Survey, error)
	FindBySlug(ctx context.Context, slug string) (*entity.Survey, error)
	Update(ctx context.Context, survey *entity.Survey) error
	Delete(ctx context.Context, id string) error

	// New methods for pagination and filtering
	List(ctx context.Context, pagi *list.PagiRequest, filter *entity.SurveyListRequest) (*list.PagiResponse[*entity.ViewSurveyList], error)
	CountParticipants(ctx context.Context, surveyID string) (int, error)

	// SurveyResponse operations
	ResponseCreate(ctx context.Context, response *entity.SurveyResponse) error
	ResponseFindByID(ctx context.Context, id string) (*entity.SurveyResponse, error)
	ResponseUpdate(ctx context.Context, response *entity.SurveyResponse) error
	ResponseList(ctx context.Context, surveyID string) ([]*entity.SurveyResponse, error)
	ResponseStats(ctx context.Context, surveyID string) (map[string]interface{}, error)
}
