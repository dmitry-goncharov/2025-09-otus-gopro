package hw09structvalidator

import (
	"fmt"
	"reflect"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}
	msg := fmt.Sprintln("Validation failed")
	for _, err := range v {
		msg += fmt.Sprintln("Field", err.Field, "Err", err.Err.Error())
	}
	return msg
}

const (
	ValidateTag = "validate"
)

func Validate(v interface{}) error {
	refType := reflect.TypeOf(v)

	if refType.Kind() != reflect.Struct {
		return nil
	}

	errors := ValidationErrors{}

	refVal := reflect.ValueOf(v)
	n := refVal.NumField()
	for i := range n {
		fldType := refType.Field(i)
		fldValue := refVal.Field(i)

		tagValue := fldType.Tag.Get(ValidateTag)

		if len(tagValue) == 0 {
			continue
		}

		errors = append(errors, validateAny(tagValue, fldType.Name, fldValue)...)
	}

	if len(errors) == 0 {
		return nil
	}

	return errors
}
