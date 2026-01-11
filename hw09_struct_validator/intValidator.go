package hw09structvalidator

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

const (
	IntRuleMin      = "min"
	IntRuleMax      = "max"
	IntRuleIn       = "in"
	IntRuleInValSep = ","
)

type IntRuleMinValidator struct {
	fldName string
	fldVal  int
	num     int
}

func (v *IntRuleMinValidator) Validate() *ValidationError {
	if v.fldVal < v.num {
		return &ValidationError{
			Field: v.fldName,
			Err:   fmt.Errorf("%d should not be less than %d", v.fldVal, v.num),
		}
	}
	return nil
}

type IntRuleMaxValidator struct {
	fldName string
	fldVal  int
	num     int
}

func (v *IntRuleMaxValidator) Validate() *ValidationError {
	if v.fldVal > v.num {
		return &ValidationError{
			Field: v.fldName,
			Err:   fmt.Errorf("%d should not be more than %d", v.fldVal, v.num),
		}
	}
	return nil
}

type IntRuleInValidator struct {
	fldName string
	fldVal  int
	nums    []int
	parts   []string
}

func (v *IntRuleInValidator) Validate() *ValidationError {
	if !slices.Contains(v.nums, v.fldVal) {
		return &ValidationError{
			Field: v.fldName,
			Err:   fmt.Errorf("%d should be in %s", v.fldVal, v.parts),
		}
	}
	return nil
}

func intValidator(tagName, tagVal, fldName string, fldVal int) (Validator, error) {
	switch tagName {
	case IntRuleMin:
		return intRuleMinValidator(tagVal, fldName, fldVal)
	case IntRuleMax:
		return intRuleMaxValidator(tagVal, fldName, fldVal)
	case IntRuleIn:
		return intRuleInValidator(tagVal, fldName, fldVal)
	}
	return nil, fmt.Errorf("invalid IntRule tag name: %q of field: %q", tagName, fldName)
}

func intRuleMinValidator(tagVal, fldName string, fldVal int) (Validator, error) {
	num, err := strconv.Atoi(tagVal)
	if err != nil {
		return nil, fmt.Errorf("invalid IntRuleMin tag value: %q of field: %q, err: %w", tagVal, fldName, err)
	}
	return &IntRuleMinValidator{
		fldName: fldName,
		fldVal:  fldVal,
		num:     num,
	}, nil
}

func intRuleMaxValidator(tagVal, fldName string, fldVal int) (Validator, error) {
	num, err := strconv.Atoi(tagVal)
	if err != nil {
		return nil, fmt.Errorf("invalid IntRuleMax tag value: %q of field: %q, err: %w", tagVal, fldName, err)
	}
	return &IntRuleMaxValidator{
		fldName: fldName,
		fldVal:  fldVal,
		num:     num,
	}, nil
}

func intRuleInValidator(tagVal, fldName string, fldVal int) (Validator, error) {
	parts := strings.Split(tagVal, IntRuleInValSep)
	nums := make([]int, 0, len(parts))
	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid IntRuleIn tag value: %q of field: %q, err: %w", tagVal, fldName, err)
		}
		nums = append(nums, num)
	}
	return &IntRuleInValidator{
		fldName: fldName,
		fldVal:  fldVal,
		nums:    nums,
		parts:   parts,
	}, nil
}
