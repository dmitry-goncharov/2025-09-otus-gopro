package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	ValidateTag    = "validate"
	RulesSeparator = "|"
	RuleSeparator  = ":"
)

func buildValidators(v interface{}) (Validators, error) {
	refType := reflect.TypeOf(v)

	if refType.Kind() != reflect.Struct {
		return nil, nil
	}

	validators := Validators{}

	refVal := reflect.ValueOf(v)
	n := refVal.NumField()
	for i := range n {
		fldType := refType.Field(i)
		fldValue := refVal.Field(i)

		tagValue := fldType.Tag.Get(ValidateTag)

		if len(tagValue) == 0 {
			continue
		}

		anyvalidators, err := anyValidators(tagValue, fldType.Name, fldValue)
		if err != nil {
			return nil, err
		}

		validators = append(validators, anyvalidators...)
	}

	return validators, nil
}

func anyValidators(tagValue, fldName string, fldValue reflect.Value) (Validators, error) {
	validators := Validators{}

	parts := strings.Split(tagValue, RulesSeparator)
	for _, tagVal := range parts {
		switch fldValue.Kind() { //nolint:exhaustive
		case reflect.Slice:
			n := fldValue.Len()
			for i := range n {
				anyvalidators, err := anyValidators(tagVal, fldName, fldValue.Index(i))
				if err != nil {
					return nil, err
				}
				validators = append(validators, anyvalidators...)
			}
		default:
			validator, err := primeValidator(tagVal, fldName, fldValue)
			if err != nil {
				return nil, err
			}
			validators = append(validators, validator)
		}
	}

	return validators, nil
}

func primeValidator(tagValue, fldName string, fldValue reflect.Value) (Validator, error) {
	parts := strings.Split(tagValue, RuleSeparator)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid rule tag value: %q of field: %q", tagValue, fldName)
	}
	tagName := parts[0]
	tagVal := parts[1]
	switch fldValue.Kind() { //nolint:exhaustive
	case reflect.Int:
		return intValidator(tagName, tagVal, fldName, int(fldValue.Int()))
	case reflect.String:
		return stringValidator(tagName, tagVal, fldName, fldValue.String())
	}
	return nil, nil
}
