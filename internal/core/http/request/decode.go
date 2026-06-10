package decode

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/dirijable/coworking-api/internal/core/error/apperror"
	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New(validator.WithRequiredStructEnabled())
)

func ValidateAndDecode(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("JSON decode: %w", err)
	}
	if err := validate.Struct(dest); err != nil {
		if validationErrors, ok := errors.AsType[validator.ValidationErrors](err); ok {
			return extractValidationErrors(&validationErrors)
		}
		return fmt.Errorf("validation error: %w", err)
	}
	return nil
}

func extractValidationErrors(validationErrors *validator.ValidationErrors) error {
	errs := make(map[string]string, len(*validationErrors))
	for _, ve := range *validationErrors {
		errs[ve.Field()] = ve.Tag()
	}
	return fmt.Errorf("%w: %w", apperror.ErrBadRequest, apperror.ValidationError{Fields: errs})
}
