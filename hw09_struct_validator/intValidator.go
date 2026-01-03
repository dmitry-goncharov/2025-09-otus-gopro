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

func validateInt(tagName, tagVal, fldName string, fldVal int) *ValidationError {
	switch tagName {
	case IntRuleMin:
		return validateIntRuleMin(tagVal, fldName, fldVal)
	case IntRuleMax:
		return validateIntRuleMax(tagVal, fldName, fldVal)
	case IntRuleIn:
		return validateIntRuleIn(tagVal, fldName, fldVal)
	}
	return &ValidationError{
		Field: fldName,
		Err:   fmt.Errorf("invalid IntRule tag name %s", tagName),
	}
}

func validateIntRuleMin(tagVal, fldName string, fldVal int) *ValidationError {
	num, err := strconv.Atoi(tagVal)
	if err != nil {
		return &ValidationError{
			Field: fldName,
			Err:   fmt.Errorf("invalid IntRuleMin tag value %s", tagVal),
		}
	}
	if fldVal < num {
		return &ValidationError{
			Field: fldName,
			Err:   fmt.Errorf("%d should not be less than %d", fldVal, num),
		}
	}
	return nil
}

func validateIntRuleMax(tagVal, fldName string, fldVal int) *ValidationError {
	num, err := strconv.Atoi(tagVal)
	if err != nil {
		return &ValidationError{
			Field: fldName,
			Err:   fmt.Errorf("invalid IntRuleMax tag value %s", tagVal),
		}
	}
	if fldVal > num {
		return &ValidationError{
			Field: fldName,
			Err:   fmt.Errorf("%d should not be more than %d", fldVal, num),
		}
	}
	return nil
}

func validateIntRuleIn(tagVal, fldName string, fldVal int) *ValidationError {
	parts := strings.Split(tagVal, IntRuleInValSep)
	nums := make([]int, 0, len(parts))
	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			return &ValidationError{
				Field: fldName,
				Err:   fmt.Errorf("invalid IntRuleIn tag value %s", tagVal),
			}
		}
		nums = append(nums, num)
	}
	if !slices.Contains(nums, fldVal) {
		return &ValidationError{
			Field: fldName,
			Err:   fmt.Errorf("%d should be in %s", fldVal, parts),
		}
	}
	return nil
}
