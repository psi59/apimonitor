package models

type DefaultValidateChecker struct {
	isValidated bool `json:"-"`
}

func (c *DefaultValidateChecker) SetValidated() {
	c.isValidated = true
}

func (c *DefaultValidateChecker) SetInvalidated() {
	c.isValidated = false
}

func (c DefaultValidateChecker) IsValidated() bool {
	return c.isValidated
}

var ValidatedDefaultValidateChecker = DefaultValidateChecker{isValidated: true}
