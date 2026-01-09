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

type StringRuleLenValidator struct {
	fldName string
	fldVal  string
	strlen  int
}

func (v *StringRuleLenValidator) Validate() *ValidationError {
	if len(v.fldVal) != v.strlen {
		return &ValidationError{
			Field: v.fldName,
			Err:   fmt.Errorf("%s should be length %d", v.fldVal, v.strlen),
		}
	}
	return nil
}

type StringRuleRegexpValidator struct {
	fldName string
	fldVal  string
	ok      bool
	pattern string
}

func (v *StringRuleRegexpValidator) Validate() *ValidationError {
	if !v.ok {
		return &ValidationError{
			Field: v.fldName,
			Err:   fmt.Errorf("%s should be match with regexp pattern %s", v.fldVal, v.pattern),
		}
	}
	return nil
}

type StringRuleInValidator struct {
	fldName string
	fldVal  string
	parts   []string
}

func (v *StringRuleInValidator) Validate() *ValidationError {
	if !slices.Contains(v.parts, v.fldVal) {
		return &ValidationError{
			Field: v.fldName,
			Err:   fmt.Errorf("%s should be in %s", v.fldVal, v.parts),
		}
	}
	return nil
}

func stringValidator(tagName, tagVal, fldName, fldVal string) (Validator, error) {
	switch tagName {
	case StringRuleLen:
		return stringRuleLenValidator(tagVal, fldName, fldVal)
	case StringRuleRegexp:
		return stringRuleRegexpValidator(tagVal, fldName, fldVal)
	case StringRuleIn:
		return stringRuleInValidator(tagVal, fldName, fldVal)
	}
	return nil, fmt.Errorf("invalid StringRule tag name: %q of field: %q", tagName, fldName)
}

func stringRuleLenValidator(tagVal, fldName, fldVal string) (Validator, error) {
	strlen, err := strconv.Atoi(tagVal)
	if err != nil {
		return nil, fmt.Errorf("invalid StrRuleLen tag value: %q of field: %q, err: %w", tagVal, fldName, err)
	}
	return &StringRuleLenValidator{
		fldName: fldName,
		fldVal:  fldVal,
		strlen:  strlen,
	}, nil
}

func stringRuleRegexpValidator(tagVal, fldName, fldVal string) (Validator, error) {
	pattern := tagVal
	ok, err := regexp.MatchString(pattern, fldVal)
	if err != nil {
		return nil, fmt.Errorf("invalid StrRuleRegexp tag value: %q of field: %q, err: %w", tagVal, fldName, err)
	}
	return &StringRuleRegexpValidator{
		fldName: fldName,
		fldVal:  fldVal,
		ok:      ok,
		pattern: pattern,
	}, nil
}

func stringRuleInValidator(tagVal, fldName, fldVal string) (Validator, error) {
	parts := strings.Split(tagVal, StringRuleInValSep)
	return &StringRuleInValidator{
		fldName: fldName,
		fldVal:  fldVal,
		parts:   parts,
	}, nil
}
