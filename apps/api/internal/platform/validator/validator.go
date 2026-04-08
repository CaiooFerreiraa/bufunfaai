package validator

import (
	"errors"
	"fmt"
	"strings"

	playgroundvalidator "github.com/go-playground/validator/v10"

	sharederrors "github.com/bufunfaai/bufunfaai/apps/api/internal/shared/errors"
)

type Validator struct {
	engine *playgroundvalidator.Validate
}

func New() *Validator {
	return &Validator{
		engine: playgroundvalidator.New(playgroundvalidator.WithRequiredStructEnabled()),
	}
}

func (validator *Validator) Validate(input any) *sharederrors.AppError {
	if err := validator.engine.Struct(input); err != nil {
		var validationErrors playgroundvalidator.ValidationErrors
		if errors.As(err, &validationErrors) && len(validationErrors) > 0 {
			first := validationErrors[0]
			fieldName := strings.ToLower(first.Field())
			return sharederrors.New(
				"VALIDATION_ERROR",
				fmt.Sprintf("campo invalido: %s", fieldName),
				400,
			)
		}

		return sharederrors.New("VALIDATION_ERROR", "payload invalido", 400)
	}

	return nil
}
