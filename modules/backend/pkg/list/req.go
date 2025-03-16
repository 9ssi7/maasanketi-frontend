package list

// PagiRequest is a struct for pagination request
type PagiRequest struct {

	// Page is a pointer to int
	// Page is a query parameter for page number
	// Page is validated to be greater than 0
	Page *uint64 `query:"page" validate:"omitempty,gt=0"`

	// Limit is a pointer to int
	// Limit is a query parameter for limit number of items per page
	// Limit is validated to be greater than 0
	Limit *uint64 `query:"limit" validate:"omitempty,gt=0"`
}

// Default is a method to set default value for PagiRequest
// Default set Page to 1 if Page is nil or less than 1
// Default set Limit to 10 if Limit is nil or less than 1
func (r *PagiRequest) Default() {
	if r.Page == nil || *r.Page <= 0 {
		r.Page = new(uint64)
		*r.Page = 1
	}
	if r.Limit == nil || *r.Limit <= 0 {
		r.Limit = new(uint64)
		*r.Limit = 10
	}
}

// Offset is a method to calculate offset for pagination
func (r *PagiRequest) Offset() uint64 {
	return (*r.Page - 1) * *r.Limit
}

func (r *PagiRequest) Offset64() int64 {
	return int64((*r.Page - 1) * *r.Limit)
}

func (r *PagiRequest) Limit64() int64 {
	return int64(*r.Limit)
}
