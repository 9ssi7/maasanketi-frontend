package response

import "time"

type ViewDefault struct {
	ID        string    `json:"id"`
	SurveyID  string    `json:"surveyId"`
	StartedAt time.Time `json:"startedAt"`
	Answers   Answers   `json:"answers"`
} // @name response.ViewDefault

func (r *Response) ToView() *ViewDefault {
	return &ViewDefault{
		ID:        r.ID,
		SurveyID:  r.SurveyID,
		StartedAt: r.StartedAt,
		Answers:   r.Answers,
	}
}
