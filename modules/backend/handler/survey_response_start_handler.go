package handler

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mstrYoda/maasanketi.co/entity"
	"github.com/mstrYoda/maasanketi.co/pkg/rescode"
	"github.com/mstrYoda/maasanketi.co/repository"
)

type SurveyResponseStartRequest struct {
	SurveySlug  string `params:"slug" validate:"required,slug"`
	IsAnonymous *bool  `json:"isAnonymous" validate:"required"`
} // @name SurveyResponseStartRequest

// SurveyResponseStart starts a new survey response
// @Summary Start a new survey response
// @Description Start a new response for a survey
// @Tags survey-responses
// @Accept json
// @Produce json
// @Param slug path string true "Survey Slug"
// @Param isAnonymous body bool false "Is Anonymous"
// @Success 201 {object} entity.ViewSurveyResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /surveys/{slug}/responses [post]
func SurveyResponseStart(repo repository.SurveyRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req SurveyResponseStartRequest
		if err := c.BodyParser(&req); err != nil {
			return rescode.ValidationFailed(err)
		}
		if err := c.ParamsParser(&req); err != nil {
			return rescode.ValidationFailed(err)
		}
		survey, err := repo.FindBySlug(c.UserContext(), req.SurveySlug)
		if err != nil {
			return rescode.SurveyNotFound(errors.New("survey not found"))
		}

		// Check if survey has expired
		if survey.IsExpired() {
			return rescode.ValidationFailed(errors.New("survey has expired"))
		}

		// Create a new response
		response := &entity.SurveyResponse{
			ID:          uuid.New().String(),
			SurveyID:    survey.ID,
			Answers:     entity.Answers{},
			StartedAt:   time.Now(),
			IPAddress:   c.IP(),
			UserAgent:   c.Get("User-Agent"),
			IsCompleted: false,
			IsAnonymous: *req.IsAnonymous,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := repo.ResponseCreate(c.UserContext(), response); err != nil {
			return err
		}
		return c.Status(fiber.StatusCreated).JSON(response.ToViewSurveyResponse())
	}
}
