package graph

import "time"

type Graph struct {
	ID          string       `json:"id"`
	SurveyID    string       `json:"surveyId"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Kind        Kind         `json:"kind"`
	Content     GraphContent `json:"content"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

type GraphContent struct {
	KeyFields   []string `json:"keyFields"`
	KeyLabels   []string `json:"keyLabels"`
	ValueFields []string `json:"valueFields"`
	ValueLabels []string `json:"valueLabels"`
}

type Kind string

const (
	KindBar     Kind = "bar"
	KindLine    Kind = "line"
	KindRadar   Kind = "radar"
	KindArea    Kind = "area"
	KindScatter Kind = "scatter"
	KindPie     Kind = "pie"
)
