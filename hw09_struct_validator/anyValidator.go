package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	RulesSeparator = "|"
	RuleSeparator  = ":"
)

func validateAny(tagValue, fldName string, fldValue reflect.Value) ValidationErrors {
	errors := ValidationErrors{}

	parts := strings.Split(tagValue, RulesSeparator)
	for _, tagVal := range parts {
		switch fldValue.Kind() { //nolint:exhaustive
		case reflect.Slice:
			n := fldValue.Len()
			for i := range n {
				errors = append(errors, validateAny(tagVal, fldName, fldValue.Index(i))...)
			}
		default:
			err := validatePrime(tagVal, fldName, fldValue)
			if err != nil {
				errors = append(errors, *err)
			}
		}
	}

	return errors
}

func validatePrime(tagValue, fldName string, fldValue reflect.Value) *ValidationError {
	parts := strings.Split(tagValue, RuleSeparator)
	if len(parts) != 2 {
		return &ValidationError{
			Field: fldName,
			Err:   fmt.Errorf("invalid rule tag value %s", tagValue),
		}
	}
	tagName := parts[0]
	tagVal := parts[1]
	switch fldValue.Kind() { //nolint:exhaustive
	case reflect.Int:
		return validateInt(tagName, tagVal, fldName, int(fldValue.Int()))
	case reflect.String:
		return validateString(tagName, tagVal, fldName, fldValue.String())
	}
	return nil
}
