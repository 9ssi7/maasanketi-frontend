package survey

import (
	"errors"
	"time"

	"github.com/mstrYoda/maasanketi.co/domain/response"
	"github.com/mstrYoda/maasanketi.co/pkg/rescode"
)

// IsExpired checks if the survey has expired
func (s *Survey) IsExpired() bool {
	if s.FinishesAt == nil {
		return false // No expiration date means it never expires
	}
	return time.Now().After(*s.FinishesAt)
}

func (s *Survey) ValidateAnswers(answers response.Answers) error {
	for _, question := range s.Questions {
		if err := s.validateByType(answers, question); err != nil {
			return err
		}
	}
	return nil
}

func (s *Survey) validateByType(answers response.Answers, question Question) error {
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

func (s *Survey) validateText(answers response.Answers, question Question) error {
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

func (s *Survey) validateNumber(answers response.Answers, question Question) error {
	_, ok := answers[question.ID].(float64)
	if !ok {
		if question.Required {
			return s.newValidationError(question.ID, "This field is required")
		}
		return nil
	}
	return nil
}

func (s *Survey) validateSelect(answers response.Answers, question Question) error {
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

func (s *Survey) validateMultiSelect(answers response.Answers, question Question) error {
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

func (s *Survey) validateBoolean(answers response.Answers, question Question) error {
	_, ok := answers[question.ID].(bool)
	if !ok {
		if question.Required {
			return s.newValidationError(question.ID, "This field is required")
		}
		return nil
	}
	return nil
}

func (s *Survey) newValidationError(field string, message string) error {
	return rescode.SurveyResponseValidationFailed(errors.New(message)).SetData(rescode.R{
		"field":     field,
		"message":   message,
		"namespace": "survey.answers",
	})
}
