package resgraph

import "time"

type ResponseGraph struct {
	ID        string                 `json:"id"`
	SurveyID  string                 `json:"surveyId"`
	GraphID   string                 `json:"graphId"`
	Content   []ResponseGraphContent `json:"content"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
}

type ResponseGraphContent struct {
	Labels []string `json:"labels"`
	Values []string `json:"values"`
}
