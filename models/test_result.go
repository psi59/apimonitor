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
	TestId       string    `json:"testId" gorm:"NOT NULL"`
	IsSuccess    bool      `json:"isSuccess"`
	StatusCode   int       `json:"statusCode"`
	Response     string    `json:"response"`
	ResponseTime int64     `json:"responseTime"`
	TestedAt     time.Time `json:"testedAt"`
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
	Page          int       `json:"page"`
	NumItem       int       `json:"numItem"`
	IsSuccess     IsSuccess `json:"isSuccess"`
	StartTestedAt time.Time `json:"startTestedAt"`
	EndTestedAt   time.Time `json:"endTestedAt"`
}

func (request *TestResultListRequest) SetZeroToDefault() {
	if request.Page == 0 {
		request.Page = 1
	}
	if request.NumItem == 0 {
		request.NumItem = 20
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
