package response

import "time"

type Answers map[string]interface{}

// SurveyResponse represents a user's response to a survey
type Response struct {
	ID          string     `json:"id"`
	SurveyID    string     `json:"surveyId"`
	UserID      string     `json:"userId,omitempty"` // Optional, for authenticated users
	Answers     Answers    `json:"answers"`          // Map of question ID to answer
	StartedAt   time.Time  `json:"startedAt"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
	IPAddress   string     `json:"ipAddress"`
	UserAgent   string     `json:"userAgent"`
	IsCompleted bool       `json:"isCompleted"`
	IsAnonymous bool       `json:"isAnonymous"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
} // @name SurveyResponse

// Complete marks the survey response as completed
func (r *Response) Complete() {
	now := time.Now()
	r.CompletedAt = &now
	r.UpdatedAt = now
	r.IsCompleted = true
}
