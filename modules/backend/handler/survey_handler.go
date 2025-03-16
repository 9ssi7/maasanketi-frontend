package handler

import (
	"errors"
	"time"

	"github.com/9ssi7/slug"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mstrYoda/maasanketi.co/entity"
	"github.com/mstrYoda/maasanketi.co/pkg/list"
	"github.com/mstrYoda/maasanketi.co/pkg/rescode"
	"github.com/mstrYoda/maasanketi.co/repository"
)

// SurveyHandler handles HTTP requests related to surveys
type SurveyHandler struct {
	surveyRepo repository.SurveyRepository
}

// NewSurveyHandler creates a new SurveyHandler
func NewSurveyHandler(surveyRepo repository.SurveyRepository) *SurveyHandler {
	return &SurveyHandler{
		surveyRepo: surveyRepo,
	}
}

// RegisterRoutes registers the survey routes
func (h *SurveyHandler) RegisterRoutes(surveyGroup fiber.Router) {

	//surveyGroup.Post("/", h.CreateSurvey) Survey  creation is not allowed for now
	surveyGroup.Get("/", h.ListSurveysWithPagination)
	surveyGroup.Get("/:slug", h.GetSurvey)

	// Response routes
	surveyGroup.Post("/:slug/responses/start", h.StartSurveyResponse)
	surveyGroup.Post("/:slug/responses/:responseId/complete", h.CompleteSurveyResponse)
	surveyGroup.Put("/:slug/responses/:responseId", h.UpdateSurveyResponse)
	surveyGroup.Get("/:slug/responses", h.ListSurveyResponses)
	surveyGroup.Get("/:slug/stats", h.GetSurveyStats)
}

// CreateSurvey creates a new survey
// @Summary Create a new survey
// @Description Create a new survey with the provided data
// @Tags surveys
// @Accept json
// @Produce json
// @Param survey body entity.Survey true "Survey data"
// @Success 201 {object} entity.Survey
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /surveys [post]
func (h *SurveyHandler) CreateSurvey(c *fiber.Ctx) error {
	survey := new(entity.Survey)
	if err := c.BodyParser(survey); err != nil {
		return rescode.ValidationFailed(err)
	}

	// Generate a new UUID for the survey
	survey.ID = uuid.New().String()
	survey.CreatedAt = time.Now()
	survey.UpdatedAt = time.Now()

	// If no questions are provided, use the default survey questions
	if len(survey.Questions) == 0 {
		defaultSurvey := entity.DefaultSurvey()
		survey.Questions = defaultSurvey.Questions
		survey.Title = defaultSurvey.Title
		survey.Description = defaultSurvey.Description
		survey.MinCompletionTimeMin = defaultSurvey.MinCompletionTimeMin
		survey.CreatedBy = defaultSurvey.CreatedBy
		survey.Tags = defaultSurvey.Tags
	}

	// Initialize tags if nil
	if survey.Tags == nil {
		survey.Tags = []string{}
	}

	// Validate FinishesAt if provided
	if survey.FinishesAt != nil && survey.FinishesAt.Before(time.Now()) {
		return rescode.ValidationFailed(errors.New("finishes_at must be in the future"))
	}

	if err := h.surveyRepo.CreateSurvey(c.UserContext(), survey); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(survey)
}

// GetSurvey gets a survey by Slug
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
func (h *SurveyHandler) GetSurvey(c *fiber.Ctx) error {
	s := c.Params("slug")

	if !slug.Is(s) {
		return rescode.SlugInvalid(errors.New("invalid slug format"))
	}

	survey, err := h.surveyRepo.GetSurveyBySlug(c.UserContext(), s)
	if err != nil {
		return rescode.SurveyNotFound(errors.New("survey not found"))
	}

	return c.JSON(survey)
}

type SurveySwaggerListResponse struct {
	Page  uint64                  `json:"page"`
	Limit uint64                  `json:"limit"`
	List  []entity.SurveyListItem `json:"list"`
}

// ListSurveysWithPagination lists all surveys with pagination and filtering
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
// @Success 200 {object} SurveySwaggerListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /surveys [get]
func (h *SurveyHandler) ListSurveysWithPagination(c *fiber.Ctx) error {
	var pagi list.PagiRequest
	if err := c.QueryParser(&pagi); err != nil {
		return rescode.ValidationFailed(err)
	}
	var filter entity.SurveyListRequest
	if err := c.QueryParser(&filter); err != nil {
		return rescode.ValidationFailed(err)
	}

	response, err := h.surveyRepo.ListSurveysWithPagination(c.UserContext(), &pagi, &filter)
	if err != nil {
		return err
	}

	return c.JSON(response)
}

// ListSurveys lists all surveys (deprecated, use ListSurveysWithPagination instead)
// @Summary List all surveys
// @Description List all available surveys
// @Tags surveys
// @Accept json
// @Produce json
// @Success 200 {array} entity.Survey
// @Failure 500 {object} map[string]interface{}
// @Router /surveys [get]
func (h *SurveyHandler) ListSurveys(c *fiber.Ctx) error {
	surveys, err := h.surveyRepo.ListSurveys(c.UserContext())
	if err != nil {
		return err
	}

	return c.JSON(surveys)
}

// StartSurveyResponse starts a new survey response
// @Summary Start a new survey response
// @Description Start a new response for a survey
// @Tags survey-responses
// @Accept json
// @Produce json
// @Param slug path string true "Survey Slug"
// @Success 201 {object} entity.SurveyResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /surveys/{slug}/responses/start [post]
func (h *SurveyHandler) StartSurveyResponse(c *fiber.Ctx) error {
	surveySlug := c.Params("slug")

	// Validate slug format
	if !slug.Is(surveySlug) {
		return rescode.SlugInvalid(errors.New("invalid slug format"))
	}

	// Check if survey exists
	survey, err := h.surveyRepo.GetSurveyBySlug(c.Context(), surveySlug)
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
		SurveySlug:  surveySlug,
		Answers:     make(map[string]interface{}),
		StartedAt:   time.Now(),
		IPAddress:   c.IP(),
		UserAgent:   c.Get("User-Agent"),
		IsCompleted: false,
		IsAnonymous: true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.surveyRepo.CreateSurveyResponse(c.Context(), response); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// UpdateSurveyResponse updates a survey response
// @Summary Update a survey response
// @Description Update a survey response with the provided data
// @Tags survey-responses
// @Accept json
// @Produce json
// @Param slug path string true "Survey Slug"
// @Param responseId path string true "Response ID"
// @Param response body map[string]interface{} true "Response data"
// @Success 200 {object} entity.SurveyResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /surveys/{slug}/responses/{responseId} [put]
func (h *SurveyHandler) UpdateSurveyResponse(c *fiber.Ctx) error {
	surveySlug := c.Params("slug")
	responseID := c.Params("responseId")

	// Validate UUID format
	if !slug.Is(surveySlug) {
		return rescode.SlugInvalid(errors.New("invalid slug format"))
	}

	if _, err := uuid.Parse(responseID); err != nil {
		return rescode.IDInvalid(errors.New("invalid response UUID format"))
	}

	// Get the survey
	survey, err := h.surveyRepo.GetSurveyBySlug(c.Context(), surveySlug)
	if err != nil {
		return rescode.SurveyNotFound(errors.New("survey not found"))
	}

	// Check if survey has expired
	if survey.IsExpired() {
		return rescode.ValidationFailed(errors.New("survey has expired"))
	}

	// Get the existing response
	response, err := h.surveyRepo.GetSurveyResponseByID(c.Context(), responseID)
	if err != nil {
		return rescode.SurveyNotFound(errors.New("survey response not found"))
	}

	// Ensure the response belongs to the specified survey
	if response.SurveySlug != surveySlug {
		return rescode.SurveyNotFound(errors.New("survey response does not belong to the specified survey"))
	}

	// Parse the answers from the request body
	var answers map[string]interface{}
	if err := c.BodyParser(&answers); err != nil {
		return rescode.ValidationFailed(err)
	}

	// Update the response
	response.Answers = answers
	response.UpdatedAt = time.Now()

	if err := h.surveyRepo.UpdateSurveyResponse(c.Context(), response); err != nil {
		return err
	}

	return c.JSON(response)
}

// CompleteSurveyResponse completes a survey response
// @Summary Complete a survey response
// @Description Mark a survey response as completed
// @Tags survey-responses
// @Accept json
// @Produce json
// @Param slug path string true "Survey Slug"
// @Param responseId path string true "Response ID"
// @Success 200 {object} entity.SurveyResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /surveys/{slug}/responses/{responseId}/complete [post]
func (h *SurveyHandler) CompleteSurveyResponse(c *fiber.Ctx) error {
	surveySlug := c.Params("slug")
	responseID := c.Params("responseId")

	// Validate UUID format
	if !slug.Is(surveySlug) {
		return rescode.SlugInvalid(errors.New("invalid slug format"))
	}

	if _, err := uuid.Parse(responseID); err != nil {
		return rescode.IDInvalid(errors.New("invalid response UUID format"))
	}

	// Get the survey
	survey, err := h.surveyRepo.GetSurveyBySlug(c.Context(), surveySlug)
	if err != nil {
		return rescode.SurveyNotFound(errors.New("survey not found"))
	}

	// Check if survey has expired
	if survey.IsExpired() {
		return rescode.ValidationFailed(errors.New("survey has expired"))
	}

	// Get the response
	response, err := h.surveyRepo.GetSurveyResponseByID(c.Context(), responseID)
	if err != nil {
		return rescode.SurveyNotFound(errors.New("survey response not found"))
	}

	// Ensure the response belongs to the specified survey
	if response.SurveySlug != surveySlug {
		return rescode.SurveyNotFound(errors.New("survey response does not belong to the specified survey"))
	}

	// Complete the response
	if err := response.Complete(survey); err != nil {
		return err
	}

	// Update the response in the repository
	if err := h.surveyRepo.UpdateSurveyResponse(c.Context(), response); err != nil {
		return err
	}

	return c.JSON(response)
}

// ListSurveyResponses lists all responses for a survey
// @Summary List all responses for a survey
// @Description List all responses for a specific survey
// @Tags survey-responses
// @Accept json
// @Produce json
// @Param slug path string true "Survey Slug"
// @Success 200 {array} entity.SurveyResponse
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /surveys/{slug}/responses [get]
func (h *SurveyHandler) ListSurveyResponses(c *fiber.Ctx) error {
	surveySlug := c.Params("slug")

	// Validate UUID format
	if !slug.Is(surveySlug) {
		return rescode.SlugInvalid(errors.New("invalid slug format"))
	}

	// Check if survey exists
	_, err := h.surveyRepo.GetSurveyBySlug(c.UserContext(), surveySlug)
	if err != nil {
		return err
	}

	// Get all responses for the survey
	responses, err := h.surveyRepo.ListSurveyResponses(c.UserContext(), surveySlug)
	if err != nil {
		return err
	}

	return c.JSON(responses)
}

// GetSurveyStats gets statistics for a survey
// @Summary Get statistics for a survey
// @Description Get aggregated statistics for a specific survey
// @Tags survey-responses
// @Accept json
// @Produce json
// @Param slug path string true "Survey Slug"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /surveys/{slug}/stats [get]
func (h *SurveyHandler) GetSurveyStats(c *fiber.Ctx) error {
	surveySlug := c.Params("slug")

	// Validate UUID format
	if !slug.Is(surveySlug) {
		return rescode.SlugInvalid(errors.New("invalid slug format"))
	}

	// Check if survey exists
	_, err := h.surveyRepo.GetSurveyBySlug(c.UserContext(), surveySlug)
	if err != nil {
		return err
	}

	// Get statistics for the survey
	stats, err := h.surveyRepo.GetSurveyResponseStats(c.UserContext(), surveySlug)
	if err != nil {
		return err
	}

	return c.JSON(stats)
}

// GetSurveyBySlug gets a survey by slug
func (h *SurveyHandler) GetSurveyBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return rescode.SlugRequired(errors.New("slug is missed"))
	}
	survey, err := h.surveyRepo.GetSurveyBySlug(c.UserContext(), slug)
	if err != nil {
		return err
	}
	return c.JSON(survey)
}
