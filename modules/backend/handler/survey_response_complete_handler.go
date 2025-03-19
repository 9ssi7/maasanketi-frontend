package handler

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mstrYoda/maasanketi.co/pkg/rescode"
	"github.com/mstrYoda/maasanketi.co/repository"
)

type SurveyResponseCompleteRequest struct {
	Slug       string                 `params:"slug" validate:"required,slug"`
	ResponseID string                 `params:"responseId" validate:"required,uuid"`
	Answers    map[string]interface{} `json:"answers" validate:"required"`
} // @name SurveyResponseCompleteRequest

// CompleteSurveyResponse completes a survey response
// @Summary Complete a survey response
// @Description Complete a survey response with the provided data
// @Tags survey-responses
// @Accept json
// @Produce json
// @Param slug path string true "Survey Slug"
// @Param responseId path string true "Response ID"
// @Param response body SurveyResponseCompleteRequest true "Response data"
// @Success 200 {object} entity.ViewSurveyResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /surveys/{slug}/responses/{responseId} [patch]
func SurveyResponseComplete(repo repository.SurveyRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req SurveyResponseCompleteRequest
		if err := c.BodyParser(&req); err != nil {
			return rescode.ValidationFailed(err)
		}
		if err := c.ParamsParser(&req); err != nil {
			return rescode.ValidationFailed(err)
		}
		survey, err := repo.FindBySlug(c.UserContext(), req.Slug)
		if err != nil {
			return rescode.SurveyNotFound(errors.New("survey not found"))
		}
		if survey.IsExpired() {
			return rescode.SurveyExpired(errors.New("survey has expired"))
		}
		response, err := repo.ResponseFindByID(c.UserContext(), req.ResponseID)
		if err != nil {
			return rescode.SurveyNotFound(errors.New("survey response not found"))
		}
		if response.SurveyID != survey.ID {
			return rescode.SurveyResponseNotFound(errors.New("survey response does not belong to the specified survey"))
		}
		if err := survey.ValidateAnswers(req.Answers); err != nil {
			return err
		}
		if err := response.Complete(survey); err != nil {
			return err
		}
		response.Answers = req.Answers
		response.UpdatedAt = time.Now()
		if err := repo.ResponseUpdate(c.UserContext(), response); err != nil {
			return err
		}
		return c.JSON(response)
	}
}
