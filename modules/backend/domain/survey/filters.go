package survey

// SurveyListFilters represents the request parameters for listing surveys
type SurveyListFilters struct {
	Tag         string  `query:"tag" form:"tag" validate:"omitempty,dive,required"`                                                                             // Filter by tags
	Sort        string  `query:"sort" form:"sort" validate:"omitempty,oneof=created_at_asc created_at_desc most_participants finishes_at_asc finishes_at_desc"` // Sort field (created_at, participants, etc.)
	HideExpired *string `query:"hideExpired" form:"hideExpired"`                                                                                                // Hide surveys that have expired
	Search      string  `query:"search" form:"search"`                                                                                                          // Text search query for title and description
}
