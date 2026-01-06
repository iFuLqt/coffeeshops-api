package validat

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type ValidationError struct {
	Message []string
}

func (v ValidationError) Error() string {
	return "validation error"
}

func ValidateStruct(s interface{}) error {
	var errorMessage []string
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "email":
				errorMessage = append(errorMessage, "Invalid email format")
			case "required":
				errorMessage = append(errorMessage, err.Field()+" is required")
			case "min":
				errorMessage = append(errorMessage, err.Field()+" must be at least "+err.Param())
			case "eqfield":
				errorMessage = append(errorMessage, err.Field()+" must be equal to "+err.Param()+".")
			default:
				errorMessage = append(errorMessage, err.Field()+" in invalid")
			}
		}
		return ValidationError{
			Message: errorMessage,
		}
	}
	return nil
}
