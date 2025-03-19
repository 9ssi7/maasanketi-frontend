package entity

import (
	"errors"
	"time"

	"github.com/mstrYoda/maasanketi.co/pkg/rescode"
)

func (s *Survey) newValidationError(field string, message string) error {
	return rescode.SurveyResponseValidationFailed(errors.New(message)).SetData(rescode.R{
		"field":     field,
		"message":   message,
		"namespace": "survey.answers",
	})
}

func (s *Survey) ValidateAnswers(answers Answers) error {
	for _, question := range s.Questions {
		if err := s.validateByType(answers, question); err != nil {
			return err
		}
	}
	return nil
}

func (s *Survey) validateByType(answers Answers, question Question) error {
	switch question.Type {
	case QuestionTypeText:
		return s.validateText(answers, question)
	case QuestionTypeNumber:
		return s.validateNumber(answers, question)
	case QuestionTypeSelect:
		return s.validateSelect(answers, question)
	case QuestionTypeMultiSelect:
		return s.validateMultiSelect(answers, question)
	case QuestionTypeBoolean:
		return s.validateBoolean(answers, question)
	}
	return nil
}

func (s *Survey) validateText(answers Answers, question Question) error {
	val, ok := answers[question.ID].(string)
	if !ok {
		if question.Required {
			return s.newValidationError(question.ID, "This field is required")
		}
		return nil
	}
	if val == "" {
		return s.newValidationError(question.ID, "This field is mandatory")
	}
	return nil
}

func (s *Survey) validateNumber(answers Answers, question Question) error {
	_, ok := answers[question.ID].(float64)
	if !ok {
		if question.Required {
			return s.newValidationError(question.ID, "This field is required")
		}
		return nil
	}
	return nil
}

func (s *Survey) validateSelect(answers Answers, question Question) error {
	val, ok := answers[question.ID].(string)
	if !ok {
		if question.Required {
			return s.newValidationError(question.ID, "This field is required")
		}
		return nil
	}
	if val == "" {
		return s.newValidationError(question.ID, "This field is mandatory")
	}
	return nil
}

func (s *Survey) validateMultiSelect(answers Answers, question Question) error {
	val, ok := answers[question.ID].([]interface{})
	if !ok {
		if question.Required {
			return s.newValidationError(question.ID, "This field is required")
		}
		return nil
	}
	if len(val) == 0 {
		return s.newValidationError(question.ID, "This field is mandatory")
	}
	for _, v := range val {
		if _, ok := v.(string); !ok {
			return s.newValidationError(question.ID, "This field must be an array of strings")
		}
	}
	return nil
}

func (s *Survey) validateBoolean(answers Answers, question Question) error {
	_, ok := answers[question.ID].(bool)
	if !ok {
		if question.Required {
			return s.newValidationError(question.ID, "This field is required")
		}
		return nil
	}
	return nil
}

// ValidateCompletion checks if the survey response meets the minimum completion time requirement
func (sr *SurveyResponse) ValidateCompletion(survey *Survey) error {
	if sr.CompletedAt == nil {
		return rescode.SurveyResponseNotCompleted(errors.New("survey is not completed"))
	}
	minCompletionTime := time.Duration(survey.MinCompletionTimeMin) * time.Minute
	actualCompletionTime := sr.CompletedAt.Sub(sr.StartedAt)
	if actualCompletionTime < minCompletionTime {
		return rescode.SurveyCompletionTimeTooShort(errors.New("survey completed too quickly, minimum completion time not met"))
	}

	return nil
}
