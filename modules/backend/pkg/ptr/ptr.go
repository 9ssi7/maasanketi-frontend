package ptr

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Int returns a pointer to the given int.
// Example: ptr.Int(42)
// Returns: *int
func UUID(id uuid.UUID) *uuid.UUID {
	return &id
}

// String returns a pointer to the given string.
// Example: ptr.String("hello")
// Returns: *string
func String(s string) *string {
	return &s
}

// Time returns a pointer to the given time.
// Example: ptr.Time(time.Now())
// Returns: *time.Time
func Time(t time.Time) *time.Time {
	return &t
}

// Int returns a pointer to the given int.
// Example: ptr.Int(42)
// Returns: *int
func Int(i int) *int {
	return &i
}

// Int returns a pointer to the given int.
// Example: ptr.Int64(42)
// Returns: *int
func Int64(i int64) *int64 {
	return &i
}

func Float64(f float64) *float64 {
	return &f
}

// Bool returns a pointer to the given bool.
// Example: ptr.Bool(true)
// Returns: *bool
func Bool(b bool) *bool {
	return &b
}

func StringRef(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func StrToSlice(s *string) []string {
	if s == nil {
		return nil
	}
	return strings.Split(*s, ",")
}
