package validation

import (
	"reflect"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func validateUUID(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(uuid.UUID); ok {
		if valuer.String() == uuid.Nil.String() {
			return nil
		}
		return valuer.String()
	}
	return nil
}

func validateSlug(fl validator.FieldLevel) bool {
	matched, _ := regexp.MatchString(slugRegexp, fl.Field().String())
	return matched
}
