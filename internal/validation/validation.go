package validation

import (
	"fmt"
	"regexp"
)

type Rule func(s string) error

type Rules []Rule

type Validator struct {
	rules Rules
}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) Add(rule Rule) {
	v.rules = append(v.rules, rule)
}

func (v *Validator) Validate(s string) []error {
	var errors []error
	for _, rule := range v.rules {
		if err := rule(s); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

func ValidateMaxLength(maxLength int) Rule {
	return func(s string) error {
		if len(s) > maxLength {
			return fmt.Errorf("input must be less than or equal to %d characters", maxLength)
		}
		return nil
	}
}

func ValidateMinLength(minLength int) Rule {
	return func(s string) error {
		if len(s) < minLength {
			return fmt.Errorf("input must be more than or equal to %d characters", minLength)
		}
		return nil
	}
}

func ValidateRegexp(regex *regexp.Regexp) Rule {
	return func(s string) error {
		if !regex.MatchString(s) {
			return fmt.Errorf("input failed regexp")
		}
		return nil
	}
}
