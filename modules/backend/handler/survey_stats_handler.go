package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mstrYoda/maasanketi.co/pkg/rescode"
	"github.com/mstrYoda/maasanketi.co/repository"
)

type SurveyStatsRequest struct {
	Slug string `params:"slug" validate:"required,slug"`
} // @name SurveyStatsRequest

// SurveyStats gets statistics for a survey
// @Summary Get statistics for a survey
// @Description Get aggregated statistics for a specific survey
// @Tags surveys
// @Accept json
// @Produce json
// @Param slug path string true "Survey Slug"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /surveys/{slug}/stats [get]
func SurveyStats(repo repository.SurveyRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req SurveyStatsRequest
		if err := c.ParamsParser(&req); err != nil {
			return rescode.ValidationFailed(err)
		}
		survey, err := repo.FindBySlug(c.UserContext(), req.Slug)
		if err != nil {
			return err
		}
		stats, err := repo.ResponseStats(c.UserContext(), survey.ID)
		if err != nil {
			return err
		}
		return c.JSON(stats)
	}
}
