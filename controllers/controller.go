package controllers

import "github.com/go-playground/validator/v10"

type ValidationErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func ValidateStruct[T any](payload T) []*ValidationErrorResponse {
	var errors []*ValidationErrorResponse
	var validate = validator.New()

	err := validate.Struct(payload)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, &ValidationErrorResponse{
				Field: err.StructNamespace(),
				Tag:   err.Tag(),
				Value: err.Param(),
			})
		}
	}

	return errors
}
