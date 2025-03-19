package handler

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mstrYoda/maasanketi.co/domain/response"
	"github.com/mstrYoda/maasanketi.co/domain/survey"
	"github.com/mstrYoda/maasanketi.co/pkg/rescode"
)

type ResponseCompleteRepo interface {
	FindByID(ctx context.Context, id string) (*response.Response, error)
	Update(ctx context.Context, response *response.Response) error
}

type ResponseCompleteSurveyRepo interface {
	FindBySlug(ctx context.Context, slug string) (*survey.Survey, error)
}

type ResponseCompleteRequest struct {
	Slug       string                 `params:"slug" validate:"required,slug"`
	ResponseID string                 `params:"responseId" validate:"required,uuid"`
	Answers    map[string]interface{} `json:"answers" validate:"required"`
} // @name ResponseCompleteRequest

// ResponseCompleteRequest completes a survey response
// @Summary Complete a survey response
// @Description Complete a survey response with the provided data
// @Tags survey-responses
// @Accept json
// @Produce json
// @Param slug path string true "Survey Slug"
// @Param responseId path string true "Response ID"
// @Param response body ResponseCompleteRequest true "Response data"
// @Success 200 {object} response.ViewDefault
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /surveys/{slug}/responses/{responseId} [patch]
func ResponseComplete(surveyRepo ResponseCompleteSurveyRepo, responseRepo ResponseCompleteRepo) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req ResponseCompleteRequest
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
		if survey.IsExpired() {
			return rescode.SurveyExpired(errors.New("survey has expired"))
		}
		res, err := responseRepo.FindByID(c.UserContext(), req.ResponseID)
		if err != nil {
			return rescode.SurveyNotFound(errors.New("survey response not found"))
		}
		if res.SurveyID != survey.ID {
			return rescode.SurveyResponseNotFound(errors.New("survey response does not belong to the specified survey"))
		}
		if err := survey.ValidateAnswers(req.Answers); err != nil {
			return err
		}
		if !survey.IsCompletable() {
			return rescode.SurveyCompletionTimeTooShort(errors.New("survey completion time is too short"))
		}
		res.Complete()
		res.Answers = req.Answers
		res.UpdatedAt = time.Now()
		if err := responseRepo.Update(c.UserContext(), res); err != nil {
			return err
		}
		return c.JSON(res)
	}
}
