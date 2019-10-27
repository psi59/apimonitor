package models

import (
	"time"

	"github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/rsvalid"
)

type WebServiceTestResult struct {
	DefaultValidateChecker
	Id               string    `json:"id"`
	WebServiceTestId int64     `json:"web_service_test_id"`
	IsSuccess        bool      `json:"is_success"`
	StatusCode       int       `json:"status_code"`
	ResponseTime     int64     `json:"response_time"`
	TestedAt         time.Time `json:"tested_at"`
}

func (result WebServiceTestResult) Validate() error {
	if rsvalid.IsZero(result.Id, result.WebServiceTestId, result.StatusCode, result.ResponseTime, result.TestedAt) {
		return rserrors.ErrInvalidParameter
	}
	result.SetValidated()
	return nil
}
