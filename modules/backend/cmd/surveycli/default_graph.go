package main

import (
	"time"

	"github.com/mstrYoda/maasanketi.co/domain/graph"
)

func DefaultGraphs(surveyID string) []*graph.Graph {
	t := time.Now()
	return []*graph.Graph{
		{
			SurveyID:    surveyID,
			Title:       "Salary by Employment Type",
			Description: "This graph shows the salary distribution by employment type.",
			Kind:        graph.KindBar,
			Content: graph.GraphContent{
				KeyFields:   []graph.Key{{Field: "employment_type", Label: "Employment Type"}},
				ValueFields: []graph.Value{{Field: "salary", Label: "Salary"}},
			},
			CreatedAt: t,
			UpdatedAt: t,
		},
		{
			SurveyID:    surveyID,
			Title:       "Salary by Job Title",
			Description: "This graph shows the salary distribution by job title.",
			Kind:        graph.KindBar,
			Content: graph.GraphContent{
				KeyFields:   []graph.Key{{Field: "job_title", Label: "Job Title"}},
				ValueFields: []graph.Value{{Field: "salary", Label: "Salary"}},
			},
			CreatedAt: t,
			UpdatedAt: t,
		},
	}
}
