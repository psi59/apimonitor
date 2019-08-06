package rsmodel

import (
	"math"

	jsoniter "github.com/json-iterator/go"
)

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

type PaginatedList struct {
	CurrentPage int
	NumItem     int
	TotalCount  int
	Items       interface{}
}

func (list PaginatedList) MarshalJSON() ([]byte, error) {
	m := struct {
		TotalCount  int         `json:"total_count"`
		TotalPage   int         `json:"total_page"`
		CurrentPage int         `json:"current_page"`
		HasNextPage bool        `json:"has_next_page"`
		Items       interface{} `json:"items"`
	}{
		TotalCount:  list.TotalCount,
		CurrentPage: list.CurrentPage,
		Items:       list.Items,
	}

	totalPageBeforeCeil := list.TotalCount / list.NumItem
	m.TotalPage = int(math.Ceil(float64(totalPageBeforeCeil)))
	if m.TotalPage == 0 {
		m.TotalPage = 1
	}
	m.HasNextPage = list.CurrentPage < m.TotalPage

	return jsoniter.Marshal(m)
}
