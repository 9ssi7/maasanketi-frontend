package entity

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

// Conditional represents a condition for showing a question
type Conditional struct {
	QuestionID string `json:"questionId"`
	Value      string `json:"value"`
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

type Answers map[string]interface{}

// SurveyResponse represents a user's response to a survey
type SurveyResponse struct {
	ID          string     `json:"id"`
	SurveyID    string     `json:"surveyId"`
	UserID      string     `json:"userId,omitempty"` // Optional, for authenticated users
	Answers     Answers    `json:"answers"`          // Map of question ID to answer
	StartedAt   time.Time  `json:"startedAt"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
	IPAddress   string     `json:"ipAddress"`
	UserAgent   string     `json:"userAgent"`
	IsCompleted bool       `json:"isCompleted"`
	IsAnonymous bool       `json:"isAnonymous"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
} // @name SurveyResponse

// GenerateSlug generates a slug from the survey title
func (s *Survey) GenerateSlug() {
	s.Slug = slug.New(s.Title+"-"+s.ID[:8], slug.TR)
}

// DefaultSurvey returns a default survey with predefined questions
func DefaultSurvey() *Survey {
	createdBy := "MaaşAnketi.co"
	// Set finishes_at to 3 months from now
	finishesAt := time.Now().AddDate(0, 3, 0)
	return &Survey{
		Title:                "Software Engineer Salary Survey",
		Description:          "This survey collects anonymous data about software engineer salaries and working conditions.",
		MinCompletionTimeMin: 3, // Minimum 3 minutes to complete the survey
		CreatedBy:            &createdBy,
		Tags:                 []string{"software", "engineering", "salary"},
		FinishesAt:           &finishesAt,
		Questions: []Question{
			{
				ID:       "employment_type",
				Text:     "What is your employment type?",
				Type:     QuestionTypeSelect,
				Required: true,
				Order:    1,
				Options: []Option{
					{ID: "full_time", Value: "Full-time employee"},
					{ID: "part_time", Value: "Part-time employee"},
					{ID: "freelancer", Value: "Freelancer/Contractor"},
					{ID: "other", Value: "Other"},
				},
			},
			{
				ID:          "salary",
				Text:        "What is your annual salary in your local currency?",
				Type:        QuestionTypeNumber,
				Required:    true,
				Order:       2,
				Placeholder: "e.g., 100000",
				Conditional: &Conditional{
					QuestionID: "employment_type",
					Value:      "full_time,part_time",
				},
			},
			{
				ID:          "freelance_rate",
				Text:        "What is your hourly rate in your local currency?",
				Type:        QuestionTypeNumber,
				Required:    true,
				Order:       3,
				Placeholder: "e.g., 100",
				Conditional: &Conditional{
					QuestionID: "employment_type",
					Value:      "freelancer",
				},
			},
			{
				ID:       "currency",
				Text:     "What currency are you paid in?",
				Type:     QuestionTypeSelect,
				Required: true,
				Order:    4,
				Options: []Option{
					{ID: "usd", Value: "USD"},
					{ID: "eur", Value: "EUR"},
					{ID: "gbp", Value: "GBP"},
					{ID: "try", Value: "TRY"},
					// Add more currencies as needed
				},
			},
			{
				ID:          "job_title",
				Text:        "What is your job title?",
				Type:        QuestionTypeText,
				Required:    true,
				Order:       5,
				Placeholder: "e.g., Senior Software Engineer",
			},
			{
				ID:          "total_experience",
				Text:        "How many years of total professional experience do you have?",
				Type:        QuestionTypeNumber,
				Required:    true,
				Order:       6,
				Placeholder: "e.g., 5",
			},
			{
				ID:          "company_experience",
				Text:        "How many years have you been with your current company?",
				Type:        QuestionTypeNumber,
				Required:    true,
				Order:       7,
				Placeholder: "e.g., 2",
			},
			{
				ID:          "it_department_size",
				Text:        "How many employees are in your IT/Engineering department?",
				Type:        QuestionTypeSelect,
				Required:    true,
				Order:       8,
				Placeholder: "Select a range",
				Options: []Option{
					{ID: "1-10", Value: "1-10"},
					{ID: "11-50", Value: "11-50"},
					{ID: "51-200", Value: "51-200"},
					{ID: "201-1000", Value: "201-1000"},
					{ID: "1001+", Value: "1001+"},
					{ID: "unknown", Value: "I don't know"},
				},
			},
			{
				ID:          "company_size",
				Text:        "How many employees are in your company overall?",
				Type:        QuestionTypeSelect,
				Required:    true,
				Order:       9,
				Placeholder: "Select a range",
				Options: []Option{
					{ID: "1-10", Value: "1-10"},
					{ID: "11-50", Value: "11-50"},
					{ID: "51-200", Value: "51-200"},
					{ID: "201-1000", Value: "201-1000"},
					{ID: "1001-2500", Value: "1001-2500"},
					{ID: "2501+", Value: "2501+"},
					{ID: "unknown", Value: "I don't know"},
				},
			},
			{
				ID:          "company_type",
				Text:        "What type of company do you work for?",
				Type:        QuestionTypeSelect,
				Required:    true,
				Order:       10,
				Placeholder: "Select company type",
				Options: []Option{
					{ID: "corporate", Value: "Corporate/Enterprise"},
					{ID: "startup", Value: "Startup"},
					{ID: "small_business", Value: "Small Business"},
					{ID: "family_owned", Value: "Family-owned Business"},
					{ID: "government", Value: "Government/Public Sector"},
					{ID: "nonprofit", Value: "Non-profit"},
					{ID: "other", Value: "Other"},
				},
			},
			{
				ID:          "company_sector",
				Text:        "What sector does your company operate in?",
				Type:        QuestionTypeSelect,
				Required:    true,
				Order:       11,
				Placeholder: "Select sector",
				Options: []Option{
					{ID: "tech", Value: "Technology/Software"},
					{ID: "finance", Value: "Finance/Banking"},
					{ID: "healthcare", Value: "Healthcare"},
					{ID: "education", Value: "Education"},
					{ID: "ecommerce", Value: "E-commerce/Retail"},
					{ID: "manufacturing", Value: "Manufacturing"},
					{ID: "media", Value: "Media/Entertainment"},
					{ID: "agency", Value: "Agency/Consulting"},
					{ID: "telecom", Value: "Telecommunications"},
					{ID: "other", Value: "Other"},
				},
			},
			{
				ID:          "salary_increase",
				Text:        "What was your latest salary increase rate (percentage)?",
				Type:        QuestionTypeNumber,
				Required:    false,
				Order:       12,
				Placeholder: "e.g., 10",
			},
			{
				ID:       "education",
				Text:     "What is your highest level of education?",
				Type:     QuestionTypeSelect,
				Required: true,
				Order:    13,
				Options: []Option{
					{ID: "high_school", Value: "High School"},
					{ID: "associate", Value: "Associate's Degree"},
					{ID: "bachelor", Value: "Bachelor's Degree"},
					{ID: "master", Value: "Master's Degree"},
					{ID: "phd", Value: "PhD or Doctorate"},
					{ID: "self_taught", Value: "Self-taught"},
					{ID: "bootcamp", Value: "Bootcamp"},
					{ID: "other", Value: "Other"},
				},
			},
			{
				ID:          "tech_stack",
				Text:        "What are your primary technologies/programming languages?",
				Type:        QuestionTypeMultiSelect,
				Required:    true,
				Order:       14,
				Placeholder: "Select all that apply",
				Options: []Option{
					{ID: "javascript", Value: "JavaScript"},
					{ID: "typescript", Value: "TypeScript"},
					{ID: "python", Value: "Python"},
					{ID: "java", Value: "Java"},
					{ID: "csharp", Value: "C#"},
					{ID: "cpp", Value: "C++"},
					{ID: "go", Value: "Go"},
					{ID: "rust", Value: "Rust"},
					{ID: "php", Value: "PHP"},
					{ID: "ruby", Value: "Ruby"},
					{ID: "swift", Value: "Swift"},
					{ID: "kotlin", Value: "Kotlin"},
					{ID: "scala", Value: "Scala"},
					{ID: "react", Value: "React"},
					{ID: "angular", Value: "Angular"},
					{ID: "vue", Value: "Vue.js"},
					{ID: "node", Value: "Node.js"},
					{ID: "dotnet", Value: ".NET"},
					{ID: "spring", Value: "Spring"},
					{ID: "django", Value: "Django"},
					{ID: "rails", Value: "Ruby on Rails"},
					{ID: "laravel", Value: "Laravel"},
					{ID: "aws", Value: "AWS"},
					{ID: "azure", Value: "Azure"},
					{ID: "gcp", Value: "Google Cloud"},
					{ID: "docker", Value: "Docker"},
					{ID: "kubernetes", Value: "Kubernetes"},
					{ID: "other", Value: "Other"},
				},
			},
			{
				ID:       "living_country",
				Text:     "In which country do you live?",
				Type:     QuestionTypeText,
				Required: true,
				Order:    15,
			},
			{
				ID:       "living_city",
				Text:     "In which city do you live?",
				Type:     QuestionTypeText,
				Required: true,
				Order:    16,
			},
			{
				ID:       "working_country",
				Text:     "In which country is your employer/client based?",
				Type:     QuestionTypeText,
				Required: true,
				Order:    17,
			},
			{
				ID:       "working_city",
				Text:     "In which city is your employer/client based?",
				Type:     QuestionTypeText,
				Required: true,
				Order:    18,
			},
			{
				ID:       "remote_work",
				Text:     "What is your work arrangement?",
				Type:     QuestionTypeSelect,
				Required: true,
				Order:    19,
				Options: []Option{
					{ID: "fully_remote", Value: "Fully remote"},
					{ID: "hybrid", Value: "Hybrid (partially remote)"},
					{ID: "onsite", Value: "Fully onsite"},
				},
			},
		},
	}
}

// Complete marks the survey response as completed
func (sr *SurveyResponse) Complete(survey *Survey) error {
	now := time.Now()
	sr.CompletedAt = &now
	sr.UpdatedAt = now
	if err := sr.ValidateCompletion(survey); err != nil {
		sr.CompletedAt = nil
		return err
	}
	sr.IsCompleted = true
	return nil
}

// IsExpired checks if the survey has expired
func (s *Survey) IsExpired() bool {
	if s.FinishesAt == nil {
		return false // No expiration date means it never expires
	}
	return time.Now().After(*s.FinishesAt)
}
