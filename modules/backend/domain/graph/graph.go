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

type ValueStrategy string

const (
	ValueStrategyCount ValueStrategy = "count"
	ValueStrategySum   ValueStrategy = "sum"
	ValueStrategyAvg   ValueStrategy = "avg"
)

type Key struct {
	Field string `json:"field"`
	Label string `json:"label"`
}

type Value struct {
	Field    string        `json:"field"`
	Label    string        `json:"label"`
	Strategy ValueStrategy `json:"strategy"`
}

type GraphContent struct {
	KeyFields   []Key   `json:"keyFields"`
	ValueFields []Value `json:"valueFields"`
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
