package handler

import (
	"github.com/gofiber/fiber/v2"
)

// SurveyResponse represents the survey data
type SurveyResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
}

// GetSurveys godoc
// @Summary Get all surveys
// @Description Get all available surveys
// @Tags surveys
// @Accept json
// @Produce json
// @Success 200 {array} SurveyResponse
// @Router /api/v1/surveys [get]
func GetSurveys(c *fiber.Ctx) error {
	surveys := []SurveyResponse{
		{
			ID:        "1",
			Title:     "Software Engineer Salary Survey 2023",
			CreatedAt: "2023-01-01T00:00:00Z",
		},
		{
			ID:        "2",
			Title:     "Frontend Developer Salary Survey 2023",
			CreatedAt: "2023-02-01T00:00:00Z",
		},
	}

	return c.JSON(surveys)
}

// GetSurvey godoc
// @Summary Get a survey by ID
// @Description Get a survey by its ID
// @Tags surveys
// @Accept json
// @Produce json
// @Param id path string true "Survey ID"
// @Success 200 {object} SurveyResponse
// @Failure 404 {object} map[string]string
// @Router /api/v1/surveys/{id} [get]
func GetSurvey(c *fiber.Ctx) error {
	id := c.Params("id")

	survey := SurveyResponse{
		ID:        id,
		Title:     "Software Engineer Salary Survey 2023",
		CreatedAt: "2023-01-01T00:00:00Z",
	}

	return c.JSON(survey)
}
