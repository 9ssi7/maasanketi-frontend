package handler

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mstrYoda/maasanketi.co/domain/response"
	"github.com/mstrYoda/maasanketi.co/domain/survey"
	"github.com/mstrYoda/maasanketi.co/pkg/rescode"
)

type ResponseStartSurveyRepo interface {
	FindBySlug(ctx context.Context, slug string) (*survey.Survey, error)
}

type ResponseStartRepo interface {
	Create(ctx context.Context, response *response.Response) error
}

type ResponseStartRequest struct {
	Slug        string `params:"slug" validate:"required,slug"`
	IsAnonymous *bool  `json:"isAnonymous" validate:"required"`
} // @name ResponseStartRequest

// ResponseStart starts a new survey response
// @Summary Start a new survey response
// @Description Start a new response for a survey
// @Tags survey-responses
// @Accept json
// @Produce json
// @Param slug path string true "Survey Slug"
// @Param isAnonymous body bool false "Is Anonymous"
// @Success 201 {object} response.ViewDefault
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /surveys/{slug}/responses [post]
func ResponseStart(surveyRepo ResponseStartSurveyRepo, responseRepo ResponseStartRepo) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req ResponseStartRequest
		if err := c.BodyParser(&req); err != nil {
			return rescode.ValidationFailed(err)
		}
		if err := c.ParamsParser(&req); err != nil {
			return rescode.ValidationFailed(err)
		}
		survey, err := surveyRepo.FindBySlug(c.UserContext(), req.Slug)
		if err != nil {
			return rescode.SurveyNotFound(errors.New("survey not found"))
		}

		// Check if survey has expired
		if survey.IsExpired() {
			return rescode.ValidationFailed(errors.New("survey has expired"))
		}

		// Create a new response
		response := &response.Response{
			ID:          uuid.New().String(),
			SurveyID:    survey.ID,
			Answers:     response.Answers{},
			StartedAt:   time.Now(),
			IPAddress:   c.IP(),
			UserAgent:   c.Get("User-Agent"),
			IsCompleted: false,
			IsAnonymous: *req.IsAnonymous,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := responseRepo.Create(c.UserContext(), response); err != nil {
			return err
		}
		return c.Status(fiber.StatusCreated).JSON(response.ToView())
	}
}
