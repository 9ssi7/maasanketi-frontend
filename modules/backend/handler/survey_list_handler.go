package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mstrYoda/maasanketi.co/entity"
	"github.com/mstrYoda/maasanketi.co/pkg/list"
	"github.com/mstrYoda/maasanketi.co/pkg/rescode"
	"github.com/mstrYoda/maasanketi.co/repository"
)

type SurveyListResponse struct {
	Page  uint64                  `json:"page"`
	Limit uint64                  `json:"limit"`
	List  []entity.ViewSurveyList `json:"list"`
} // @name SurveyListResponse

// SurveyList lists all surveys with pagination and filtering
// @Summary List all surveys with pagination and filtering
// @Description List all available surveys with pagination, filtering, and sorting
// @Tags surveys
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param tag query string false "Filter by tag"
// @Param sort query string false "Sort field (created_at_asc, created_at_desc, most_participants, finishes_at_asc, finishes_at_desc)"
// @Param hideExpired query bool false "Hide expired surveys"
// @Param search query string false "Text search query for title and description"
// @Success 200 {object} SurveyListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /surveys [get]
func SurveyList(repo repository.SurveyRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var pagi list.PagiRequest
		if err := c.QueryParser(&pagi); err != nil {
			return rescode.ValidationFailed(err)
		}
		var filter entity.SurveyListRequest
		if err := c.QueryParser(&filter); err != nil {
			return rescode.ValidationFailed(err)
		}

		response, err := repo.List(c.UserContext(), &pagi, &filter)
		if err != nil {
			return err
		}

		return c.JSON(response)
	}
}
