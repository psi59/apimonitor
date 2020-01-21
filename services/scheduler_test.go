package services

import (
	"net/http"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"

	"github.com/realsangil/apimonitor/models"
	"github.com/realsangil/apimonitor/pkg/rshttp"
	"github.com/realsangil/apimonitor/pkg/testutils"
)

func Test_webServiceScheduler_Run(t *testing.T) {
	testutils.SetLogConfig()
	testutils.MonkeyAll()

	ws := models.WebService{
		Id:          1,
		Host:        "realsangil.github.io",
		Schema:      "https",
		Description: "",
		Favicon:     "",
		Schedule:    models.ScheduleOneMinute,
		Tests: []models.Test{
			{
				Id:           1,
				WebServiceId: 1,
				Path:         "/",
				Method:       rshttp.MethodGet,
				ContentType:  rshttp.MIMETextHTML,
				Description:  "",
				RequestData:  nil,
				Header:       nil,
				QueryParam:   nil,
				Timeout:      0,
				Assertion: models.AssertionV1{
					StatusCode: http.StatusOK,
				},
				CreatedAt:  time.Now(),
				ModifiedAt: time.Now(),
			},
			{
				Id:           2,
				WebServiceId: 1,
				Path:         "/bio",
				Method:       rshttp.MethodGet,
				ContentType:  rshttp.MIMETextHTML,
				Description:  "",
				RequestData:  nil,
				Header:       nil,
				QueryParam:   nil,
				Timeout:      0,
				Assertion: models.AssertionV1{
					StatusCode: http.StatusNotFound,
				},
				CreatedAt:  time.Now(),
				ModifiedAt: time.Now(),
			},
			{
				Id:           2,
				WebServiceId: 1,
				Path:         "/profile/is_admin",
				Method:       rshttp.MethodGet,
				ContentType:  rshttp.MIMETextHTML,
				Description:  "",
				RequestData:  nil,
				Header:       nil,
				QueryParam:   nil,
				Timeout:      0,
				Assertion: models.AssertionV1{
					StatusCode: http.StatusNotFound,
				},
				CreatedAt:  time.Now(),
				ModifiedAt: time.Now(),
			},
		},
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	type fields struct {
		webService *models.WebService
		closeChan  chan bool
		errorChan  chan error
	}

	tests := []struct {
		name     string
		fields   fields
		testFunc func(webServiceScheduler *testScheduler) func()
		wantErr  error
	}{
		{
			name: "pass",
			fields: fields{
				webService: &ws,
				closeChan:  make(chan bool, 1),
				errorChan:  make(chan error, 1),
			},
			testFunc: func(webServiceScheduler *testScheduler) func() {
				return func() {
					_ = webServiceScheduler.Close()
				}
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schedule := &testScheduler{
				test:      tt.fields.webService,
				closeChan: tt.fields.closeChan,
			}
			time.AfterFunc(3*time.Second, tt.testFunc(schedule))
			err := schedule.Run()
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_webServiceScheduleManager_Run(t *testing.T) {
	testutils.SetLogConfig()
	testutils.MonkeyAll()
	monkey.Patch(rshttp.Do, func(request *rshttp.Request) (rshttp.Response, error) {
		return &rshttp.Response{
			StatusCode:   http.StatusOK,
			ResponseTime: 10,
			Body:         nil,
		}, nil
	})

	resultChan := make(chan *models.TestResult, 1)

	type fields struct {
		webServiceSchedulers map[interface{}]Scheduler
		resultChan           chan *models.TestResult
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "pass",
			fields: fields{
				webServiceSchedulers: map[interface{}]Scheduler{
					1: &testScheduler{
						test: &models.WebService{
							Id:       1,
							Schedule: models.ScheduleOneMinute,
							Tests: []models.Test{
								{
									Id:           1,
									WebServiceId: 1,
									Path:         "/",
									Method:       rshttp.MethodGet,
									ContentType:  rshttp.MIMEApplicationJSON,
									Timeout:      0,
									Assertion: models.AssertionV1{
										StatusCode: http.StatusOK,
									},
									CreatedAt:  time.Now(),
									ModifiedAt: time.Now(),
								},
							},
							CreatedAt:  time.Now(),
							ModifiedAt: time.Now(),
						},
						closeChan:  make(chan bool, 1),
						resultChan: resultChan,
					},
				},
				resultChan: resultChan,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &TestScheduleManager{
				testSchedulers: tt.fields.webServiceSchedulers,
				resultChan:     tt.fields.resultChan,
				closeChan:      make(chan bool, 1),
			}
			time.AfterFunc(1*time.Second, func() {
				manager.Close()
			})
			if err := manager.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
