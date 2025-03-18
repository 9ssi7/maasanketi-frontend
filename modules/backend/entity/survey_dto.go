package entity

import "time"

// ViewSurveyList represents a survey in a list view without questions
type ViewSurveyList struct {
	ID                   string     `json:"id"`
	Slug                 string     `json:"slug"`
	Title                string     `json:"title"`
	Description          string     `json:"description"`
	CreatedAt            time.Time  `json:"createdAt"`
	UpdatedAt            time.Time  `json:"updatedAt"`
	MinCompletionTimeMin int        `json:"minCompletionTimeMin"`
	CreatedBy            *string    `json:"createdBy,omitempty"`
	Tags                 []string   `json:"tags,omitempty"`
	Participants         int        `json:"participants"` // Number of completed responses
	FinishesAt           *time.Time `json:"finishesAt,omitempty"`
} // @name entity.ViewSurveyList

type ViewSurveyResponse struct {
	ID         string    `json:"id"`
	SurveySlug string    `json:"surveySlug"`
	StartedAt  time.Time `json:"startedAt"`
	Answers    Answers   `json:"answers"`
} // @name entity.ViewSurveyResponse

// SurveyListRequest represents the request parameters for listing surveys
type SurveyListRequest struct {
	Tag         string `query:"tag" form:"tag" validate:"omitempty,dive,required"`                                                                             // Filter by tags
	Sort        string `query:"sort" form:"sort" validate:"omitempty,oneof=created_at_asc created_at_desc most_participants finishes_at_asc finishes_at_desc"` // Sort field (created_at, participants, etc.)
	HideExpired bool   `query:"hideExpired" form:"hideExpired"`                                                                                                // Hide surveys that have expired
	Search      string `query:"search" form:"search"`                                                                                                          // Text search query for title and description
}

// ToListItem converts a Survey to a SurveyListItem
func (s *Survey) ToListItem(participants int) *ViewSurveyList {
	return &ViewSurveyList{
		ID:                   s.ID,
		Slug:                 s.Slug,
		Title:                s.Title,
		Description:          s.Description,
		CreatedAt:            s.CreatedAt,
		UpdatedAt:            s.UpdatedAt,
		MinCompletionTimeMin: s.MinCompletionTimeMin,
		CreatedBy:            s.CreatedBy,
		Tags:                 s.Tags,
		Participants:         participants,
		FinishesAt:           s.FinishesAt,
	}
}

func (s *SurveyResponse) ToViewSurveyResponse() *ViewSurveyResponse {
	return &ViewSurveyResponse{
		ID:         s.ID,
		SurveySlug: s.SurveySlug,
		StartedAt:  s.StartedAt,
		Answers:    s.Answers,
	}
}
