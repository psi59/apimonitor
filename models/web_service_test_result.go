package models

import (
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/rsvalid"
)

type WebServiceTestResult struct {
	DefaultValidateChecker
	Id               string    `json:"id" gorm:"Size:36"`
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

func (result WebServiceTestResult) TableName() string {
	return "web_service_test_results"
}

type WebServiceTestResultListRequest struct {
	Page          int
	NumItem       int
	IsSuccess     IsSuccess
	StartTestedAt time.Time
	EndTestedAt   time.Time
}

type IsSuccess string

func (isSuccess IsSuccess) String() string {
	return string(isSuccess)
}

func (isSuccess IsSuccess) isEmpty() bool {
	return isSuccess == ""
}

func (isSuccess IsSuccess) IsBoth() bool {
	return isSuccess.isEmpty()
}

func (isSuccess IsSuccess) Validate() error {
	if isSuccess.isEmpty() {
		return nil
	}
	if _, err := strconv.ParseBool(isSuccess.String()); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
