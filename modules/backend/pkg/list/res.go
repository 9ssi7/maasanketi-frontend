package list

// PagiResponse is a struct for pagination response
// PagiResponse contains page number, limit, total items, filtered total items, total page, and list of items
// PagiResponse is a generic struct
type PagiResponse[T any] struct {
	Page  uint64 `json:"page"`
	Limit uint64 `json:"limit"`
	List  []T    `json:"list"`
}
