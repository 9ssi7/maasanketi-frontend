package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mstrYoda/maasanketi.co/pkg/rescode"
	"github.com/mstrYoda/maasanketi.co/repository"
)

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
// @Success 200 {object} entity.Survey
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /surveys/{slug} [get]
func SurveyView(repo repository.SurveyRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req SurveyViewRequest
		if err := c.ParamsParser(&req); err != nil {
			return rescode.ValidationFailed(err)
		}
		survey, err := repo.FindBySlug(c.UserContext(), req.Slug)
		if err != nil {
			return err
		}
		return c.JSON(survey)
	}
}
