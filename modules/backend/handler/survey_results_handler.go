package handler

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/mstrYoda/maasanketi.co/domain/resgraph"
	"github.com/mstrYoda/maasanketi.co/domain/survey"
	"github.com/mstrYoda/maasanketi.co/pkg/rescode"
)

type SurveyResultsRepo interface {
	ListBySurveyID(ctx context.Context, surveyID string) ([]*resgraph.ViewList, error)
}

type SurveyResultsSurveyRepo interface {
	FindBySlug(ctx context.Context, slug string) (*survey.Survey, error)
}

type SurveyResultsRequest struct {
	Slug string `params:"slug" validate:"required,slug"`
}

// SurveyResults lists all survey results
// @Summary List all survey results
// @Description List all survey results
// @Tags survey-results
// @Accept json
// @Produce json
// @Param slug path string true "Survey Slug"
// @Success 200 {object} resgraph.ViewList
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /surveys/{slug}/results [get]
func SurveyResults(surveyRepo SurveyResultsSurveyRepo, responseGraphRepo SurveyResultsRepo) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req SurveyResultsRequest
		if err := c.ParamsParser(&req); err != nil {
			return rescode.ValidationFailed(err)
		}
		survey, err := surveyRepo.FindBySlug(c.UserContext(), req.Slug)
		if err != nil {
			return rescode.SurveyNotFound(errors.New("survey not found"))
		}
		if !survey.IsExpired() {
			return rescode.SurveyNotFinished(errors.New("survey not finished"))
		}
		responseGraphs, err := responseGraphRepo.ListBySurveyID(c.UserContext(), survey.ID)
		if err != nil {
			return err
		}
		return c.JSON(responseGraphs)
	}
}
