package main

import "github.com/mstrYoda/maasanketi.co/domain/graph"

func DefaultGraphs(surveyID string) []*graph.Graph {
	return []*graph.Graph{
		{
			SurveyID:    surveyID,
			Title:       "Salary by Employment Type",
			Description: "This graph shows the salary distribution by employment type.",
			Kind:        graph.KindBar,
			KeyFields:   []string{"employment_type"},
			ValueFields: []string{"salary"},
		},
		{
			SurveyID:    surveyID,
			Title:       "Salary by Job Title",
			Description: "This graph shows the salary distribution by job title.",
			Kind:        graph.KindBar,
			KeyFields:   []string{"job_title"},
			ValueFields: []string{"salary"},
		},
	}
}
