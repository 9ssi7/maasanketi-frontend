package survey

import (
	"time"
)

// ViewList represents a survey in a list view without questions
type ViewList struct {
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
} // @name survey.ViewList

// ToListItem converts a Survey to a ViewList
func (s *Survey) ToListView(participants int) *ViewList {
	return &ViewList{
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
