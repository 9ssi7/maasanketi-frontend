package resgraph

import "time"

type ResponseGraph struct {
	ID        string            `json:"id"`
	SurveyID  string            `json:"surveyId"`
	GraphID   string            `json:"graphId"`
	Keys      []ResponseGraphKV `json:"keys"`
	Values    []ResponseGraphKV `json:"values"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
}

type ResponseGraphKV struct {
	Field string `json:"field"`
	Value string `json:"value"`
}
