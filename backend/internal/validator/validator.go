package validator

import (
	"fmt"
	"strings"

	goValidator "github.com/go-playground/validator/v10"
)

var validate = goValidator.New()

// Validate validates a struct using go-playground/validator tags.
// Returns a map of field -> error message on failure, or nil if valid.
func Validate(s any) map[string]string {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	errs := make(map[string]string)
	for _, e := range err.(goValidator.ValidationErrors) {
		field := strings.ToLower(e.Field())
		errs[field] = fieldError(e)
	}
	return errs
}

func fieldError(e goValidator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("minimum length is %s", e.Param())
	case "max":
		return fmt.Sprintf("maximum length is %s", e.Param())
	case "oneof":
		return fmt.Sprintf("must be one of: %s", e.Param())
	default:
		return fmt.Sprintf("failed validation: %s", e.Tag())
	}
}
