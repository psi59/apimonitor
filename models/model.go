package models

type DefaultValidateChecker struct {
	isValidated bool `json:"-"`
}

func (c DefaultValidateChecker) IsValidated() bool {
	return c.isValidated
}

var ValidatedDefaultValidateChecker = DefaultValidateChecker{isValidated: true}
