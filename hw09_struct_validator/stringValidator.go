package hw09structvalidator

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const (
	StringRuleLen      = "len"
	StringRuleRegexp   = "regexp"
	StringRuleIn       = "in"
	StringRuleInValSep = ","
)

func validateString(tagName, tagVal, fldName, fldVal string) *ValidationError {
	switch tagName {
	case StringRuleLen:
		return validateStringRuleLen(tagVal, fldName, fldVal)
	case StringRuleRegexp:
		return validateStringRuleRegexp(tagVal, fldName, fldVal)
	case StringRuleIn:
		return validateStringRuleIn(tagVal, fldName, fldVal)
	}
	panic(fmt.Sprintf("invalid StringRule tag name: %q of field: %q", tagName, fldName))
}

func validateStringRuleLen(tagVal, fldName, fldVal string) *ValidationError {
	strlen, err := strconv.Atoi(tagVal)
	if err != nil {
		panic(fmt.Sprintf("invalid StrRuleLen tag value: %q of field: %q", tagVal, fldName))
	}
	if len(fldVal) != strlen {
		return &ValidationError{
			Field: fldName,
			Err:   fmt.Errorf("%s should be length %d", fldVal, strlen),
		}
	}
	return nil
}

func validateStringRuleRegexp(tagVal, fldName, fldVal string) *ValidationError {
	pattern := tagVal
	ok, err := regexp.MatchString(pattern, fldVal)
	if err != nil {
		panic(fmt.Sprintf("invalid StrRuleRegexp tag value: %q of field: %q", tagVal, fldName))
	}
	if !ok {
		return &ValidationError{
			Field: fldName,
			Err:   fmt.Errorf("%s should be match with regexp pattern %s", fldVal, pattern),
		}
	}
	return nil
}

func validateStringRuleIn(tagVal, fldName, fldVal string) *ValidationError {
	parts := strings.Split(tagVal, StringRuleInValSep)
	if !slices.Contains(parts, fldVal) {
		return &ValidationError{
			Field: fldName,
			Err:   fmt.Errorf("%s should be in %s", fldVal, parts),
		}
	}
	return nil
}
