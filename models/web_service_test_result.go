package models

import "time"

type WebServiceTestResult struct {
	DefaultValidateChecker
	Id               string    `json:"id"`
	WebServiceTestId int64     `json:"web_service_test_id"`
	IsSuccess        bool      `json:"is_success"`
	StatusCode       int       `json:"status_code"`
	ResponseTime     float64   `json:"response_time"`
	TestedAt         time.Time `json:"tested_at"`
}
