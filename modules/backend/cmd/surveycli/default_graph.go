package main

import "github.com/mstrYoda/maasanketi.co/domain/graph"

func DefaultGraphs(surveyID string) []*graph.Graph {
	return []*graph.Graph{
		{
			SurveyID:    surveyID,
			Title:       "Salary by Employment Type",
			Description: "This graph shows the salary distribution by employment type.",
			Kind:        graph.KindBar,
			Content: graph.GraphContent{
				KeyFields:   []string{"employment_type"},
				KeyLabels:   []string{"Employment Type"},
				ValueFields: []string{"salary"},
				ValueLabels: []string{"Salary"},
			},
		},
		{
			SurveyID:    surveyID,
			Title:       "Salary by Job Title",
			Description: "This graph shows the salary distribution by job title.",
			Kind:        graph.KindBar,
			Content: graph.GraphContent{
				KeyFields:   []string{"job_title"},
				KeyLabels:   []string{"Job Title"},
				ValueFields: []string{"salary"},
				ValueLabels: []string{"Salary"},
			},
		},
	}
}
