package transport

import (
	"GoFrioCalor/internal/constants"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			field := fieldError.Field()
			tag := fieldError.Tag()

			switch tag {
			case "required":
				errors[field] = fmt.Sprintf(constants.ValidationRequired, field)
			case "min":
				errors[field] = fmt.Sprintf(constants.ValidationMinLength, field, fieldError.Param())
			case "max":
				errors[field] = fmt.Sprintf(constants.ValidationMaxLength, field, fieldError.Param())
			case "len":
				errors[field] = fmt.Sprintf(constants.ValidationExactLength, field, fieldError.Param())
			case "oneof":
				errors[field] = fmt.Sprintf(constants.ValidationOneOf, field, fieldError.Param())
			case "numeric":
				errors[field] = fmt.Sprintf(constants.ValidationNumeric, field)
			case "gt":
				errors[field] = fmt.Sprintf(constants.ValidationGreaterThan, field, fieldError.Param())
			default:
				errors[field] = fmt.Sprintf(constants.ValidationInvalid, field)
			}
		}
	}

	return errors
}
