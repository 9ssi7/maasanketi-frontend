package survey

import (
	"time"

	"github.com/9ssi7/slug"
)

// QuestionType represents the type of question
type QuestionType string

const (
	// QuestionTypeText represents a free text question
	QuestionTypeText QuestionType = "text"
	// QuestionTypeNumber represents a numeric question
	QuestionTypeNumber QuestionType = "number"
	// QuestionTypeSelect represents a single select question
	QuestionTypeSelect QuestionType = "select"
	// QuestionTypeMultiSelect represents a multi-select question
	QuestionTypeMultiSelect QuestionType = "multi-select"
	// QuestionTypeBoolean represents a yes/no question
	QuestionTypeBoolean QuestionType = "boolean"
)

// Option represents a selectable option for select and multi-select questions
type Option struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

// Conditional represents a condition for showing a question
type Conditional struct {
	QuestionID string `json:"questionId"`
	Value      string `json:"value"`
}

// Question represents a survey question
type Question struct {
	ID          string       `json:"id"`
	Text        string       `json:"text"`
	Type        QuestionType `json:"type"`
	Required    bool         `json:"required"`
	Options     []Option     `json:"options,omitempty"`
	Placeholder string       `json:"placeholder,omitempty"`
	Conditional *Conditional `json:"conditional,omitempty"`
	Order       int          `json:"order"`
}

// Survey represents a survey definition
type Survey struct {
	ID                   string     `json:"id"`
	Slug                 string     `json:"slug"`
	Title                string     `json:"title"`
	Description          string     `json:"description"`
	Questions            []Question `json:"questions"`
	CreatedAt            time.Time  `json:"createdAt"`
	UpdatedAt            time.Time  `json:"updatedAt"`
	FinishesAt           *time.Time `json:"finishesAt,omitempty"`
	MinCompletionTimeMin int        `json:"minCompletionTimeMin"` // Minimum completion time in minutes
	CreatedBy            *string    `json:"createdBy,omitempty"`  // Optional creator name
	Tags                 []string   `json:"tags,omitempty"`       // Optional tags for categorization
}

// GenerateSlug generates a slug from the survey title
func (s *Survey) GenerateSlug() {
	s.Slug = slug.New(s.Title+"-"+s.ID[:8], slug.TR)
}

func (s *Survey) IsCompletable() bool {
	minCompletionTime := time.Duration(s.MinCompletionTimeMin) * time.Minute
	actualCompletionTime := time.Now().Sub(s.CreatedAt)
	return actualCompletionTime >= minCompletionTime
}
