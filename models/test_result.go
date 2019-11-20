package models

import (
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/rsmodels"
	"github.com/realsangil/apimonitor/pkg/rsvalid"
)

type TestResult struct {
	rsmodels.DefaultValidateChecker
	Id           string    `json:"id" gorm:"Size:36"`
	TestId       int64     `json:"test_id" gorm:"NOT NULL"`
	IsSuccess    bool      `json:"is_success"`
	StatusCode   int       `json:"status_code"`
	ResponseTime int64     `json:"response_time"`
	TestedAt     time.Time `json:"tested_at"`
}

func (result TestResult) Validate() error {
	if rsvalid.IsZero(result.Id, result.TestId, result.StatusCode, result.ResponseTime, result.TestedAt) {
		return rserrors.ErrInvalidParameter
	}
	result.SetValidated()
	return nil
}

func (result TestResult) TableName() string {
	return "test_results"
}

type TestResultListRequest struct {
	Page          int       `query:"page"`
	NumItem       int       `query:"num_item"`
	IsSuccess     IsSuccess `query:"is_success"`
	StartTestedAt time.Time `query:"start_tested_at"`
	EndTestedAt   time.Time `query:"end_tested_at"`
}

func (request *TestResultListRequest) SetZeroToDefault() {
	if request.Page == 0 {
		request.Page = 1
	}
	if request.NumItem == 0 {
		request.NumItem = 200
	}
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
