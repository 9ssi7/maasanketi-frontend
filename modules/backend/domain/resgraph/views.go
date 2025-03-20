package resgraph

import (
	"time"

	"github.com/mstrYoda/maasanketi.co/domain/graph"
)

type ViewList struct {
	ID           string                 `json:"id"`
	SurveyID     string                 `json:"surveyId"`
	GraphID      string                 `json:"graphId"`
	Title        string                 `json:"title"`
	Description  string                 `json:"description"`
	Kind         graph.Kind             `json:"kind"`
	GraphContent graph.GraphContent     `json:"graphContent"`
	Content      []ResponseGraphContent `json:"content"`
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt"`
}
