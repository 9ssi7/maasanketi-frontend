package validation

import (
	"context"
	"errors"
	"strings"

	"github.com/go-playground/locales/tr"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/mstrYoda/maasanketi.co/pkg/rescode"
)

type Srv struct {
	validator *validator.Validate
	uni       *ut.UniversalTranslator
}

func New() *Srv {
	v := validator.New()
	v.RegisterCustomTypeFunc(validateUUID, uuid.UUID{})
	_ = v.RegisterValidation("slug", validateSlug)
	return &Srv{validator: v, uni: ut.New(tr.New())}
}

// ValidateStruct validates the given struct.
func (s *Srv) ValidateStruct(ctx context.Context, sc interface{}) error {
	var errs []*ErrorResponse
	err := s.validator.StructCtx(ctx, sc)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			ns := s.mapStructNamespace(err.Namespace())
			if ns != "" {
				element.Namespace = ns
			}
			element.Field = err.Field()
			element.Value = err.Value()
			element.Message = err.Translate(s.getTranslator(ctx))
			errs = append(errs, &element)
		}
	}
	if len(errs) > 0 {
		return rescode.ValidationFailed(errors.New("validation failed")).SetData(errs)
	}
	return nil
}

// ValidateMap validates the giveb struct.
func (s *Srv) ValidateMap(ctx context.Context, m map[string]interface{}, rules map[string]interface{}) error {
	var errs []*ErrorResponse
	errMap := s.validator.ValidateMapCtx(ctx, m, rules)
	for key, err := range errMap {
		var element ErrorResponse
		if _err, ok := err.(validator.ValidationErrors); ok {
			for _, err := range _err {
				element.Namespace = err.Namespace()
				element.Field = err.Field()
				if element.Field == "" {
					element.Field = key
				}
				element.Value = err.Value()
				element.Message = err.Translate(s.getTranslator(ctx))
				errs = append(errs, &element)
			}
			continue
		}
	}
	if len(errs) > 0 {
		return rescode.ValidationFailed(errors.New("validation failed")).SetData(errs)
	}
	return nil
}

func (s *Srv) getTranslator(_ context.Context) ut.Translator {
	translator, found := s.uni.GetTranslator("tr")
	if !found {
		translator = s.uni.GetFallback()
	}
	return translator
}

func (s *Srv) mapStructNamespace(ns string) string {
	str := strings.Split(ns, ".")
	return strings.Join(str[1:], ".")
}
