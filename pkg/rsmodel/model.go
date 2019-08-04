package rsmodel

type Validator interface {
	Validate() error
}

type ValidateChecker interface {
	IsValidated() bool
}

type ValidatedObject interface {
	ValidateChecker
	Validator
}

var ValidatedDefaultValidateChecker = DefaultValidateChecker{isValidated: true}

type DefaultValidateChecker struct {
	isValidated bool `json:"-"`
}

func (c DefaultValidateChecker) IsValidated() bool {
	return c.isValidated
}
