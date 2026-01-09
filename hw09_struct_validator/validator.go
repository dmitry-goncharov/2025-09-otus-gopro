package hw09structvalidator

import (
	"fmt"
	"strings"
)

type Validator interface {
	Validate() *ValidationError
}

type Validators []Validator

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}
	msg := strings.Builder{}
	msg.WriteString("Validation failed")
	for _, err := range v {
		msg.WriteString(fmt.Sprintln("Field", err.Field, "Err", err.Err.Error()))
	}
	return msg.String()
}

func Validate(v interface{}) error {
	validators, err := buildValidators(v)
	if err != nil {
		return fmt.Errorf("error building validators for %v, err: %w", v, err)
	}

	errors := ValidationErrors{}
	for _, validator := range validators {
		if validator != nil {
			validationError := validator.Validate()
			if validationError != nil {
				errors = append(errors, *validationError)
			}
		}
	}

	if len(errors) == 0 {
		return nil
	}

	return errors
}
