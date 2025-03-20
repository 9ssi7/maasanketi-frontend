package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/mstrYoda/maasanketi.co/domain/survey"
	"github.com/mstrYoda/maasanketi.co/pkg/rescode"
)

type SurveyViewRepo interface {
	ViewDetail(ctx context.Context, slug string) (*survey.ViewDetail, error)
}

type SurveyViewRequest struct {
	Slug string `params:"slug" validate:"required,slug"`
} // @name SurveyViewRequest

// SurveyView gets a survey by Slug
// @Summary Get a survey by Slug
// @Description Get a survey by its Slug
// @Tags surveys
// @Accept json
// @Produce json
// @Param slug path string true "Survey Slug"
// @Success 200 {object} survey.Survey
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /surveys/{slug} [get]
func SurveyView(repo SurveyViewRepo) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req SurveyViewRequest
		if err := c.ParamsParser(&req); err != nil {
			return rescode.ValidationFailed(err)
		}
		survey, err := repo.ViewDetail(c.UserContext(), req.Slug)
		if err != nil {
			return err
		}
		return c.JSON(survey)
	}
}
