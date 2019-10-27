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
		Id:         1,
		Host:       "realsangil.github.io",
		HttpSchema: "https",
		Desc:       "",
		Favicon:    "",
		Schedule:   models.ScheduleOneMinute,
		Tests: []models.WebServiceTest{
			{
				Id:           1,
				WebServiceId: 1,
				Path:         "/",
				HttpMethod:   rshttp.MethodGet,
				ContentType:  rshttp.MIMETextHTML,
				Desc:         "",
				RequestData:  nil,
				Header:       nil,
				QueryParam:   nil,
				Timeout:      0,
				Assertion: models.AssertionV1{
					StatusCode: http.StatusOK,
				},
				Created:      time.Now(),
				LastModified: time.Now(),
			},
			{
				Id:           2,
				WebServiceId: 1,
				Path:         "/bio",
				HttpMethod:   rshttp.MethodGet,
				ContentType:  rshttp.MIMETextHTML,
				Desc:         "",
				RequestData:  nil,
				Header:       nil,
				QueryParam:   nil,
				Timeout:      0,
				Assertion: models.AssertionV1{
					StatusCode: http.StatusNotFound,
				},
				Created:      time.Now(),
				LastModified: time.Now(),
			},
			{
				Id:           2,
				WebServiceId: 1,
				Path:         "/profile/is_admin",
				HttpMethod:   rshttp.MethodGet,
				ContentType:  rshttp.MIMETextHTML,
				Desc:         "",
				RequestData:  nil,
				Header:       nil,
				QueryParam:   nil,
				Timeout:      0,
				Assertion: models.AssertionV1{
					StatusCode: http.StatusNotFound,
				},
				Created:      time.Now(),
				LastModified: time.Now(),
			},
		},
		Created:      time.Now(),
		LastModified: time.Now(),
	}

	type fields struct {
		webService *models.WebService
		closeChan  chan bool
		errorChan  chan error
	}

	tests := []struct {
		name     string
		fields   fields
		testFunc func(webServiceScheduler *webServiceScheduler) func()
		wantErr  error
	}{
		{
			name: "pass",
			fields: fields{
				webService: &ws,
				closeChan:  make(chan bool, 1),
				errorChan:  make(chan error, 1),
			},
			testFunc: func(webServiceScheduler *webServiceScheduler) func() {
				return func() {
					_ = webServiceScheduler.Close()
				}
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schedule := &webServiceScheduler{
				webService: tt.fields.webService,
				closeChan:  tt.fields.closeChan,
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
		return &rshttp.HttpResponse{
			StatusCode:   http.StatusOK,
			ResponseTime: 10,
			Body:         nil,
		}, nil
	})

	resultChan := make(chan *models.WebServiceTestResult, 1)

	type fields struct {
		webServiceSchedulers map[interface{}]WebServiceScheduler
		resultChan           chan *models.WebServiceTestResult
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "pass",
			fields: fields{
				webServiceSchedulers: map[interface{}]WebServiceScheduler{
					1: &webServiceScheduler{
						webService: &models.WebService{
							Id:       1,
							Schedule: models.ScheduleOneMinute,
							Tests: []models.WebServiceTest{
								{
									Id:           1,
									WebServiceId: 1,
									Path:         "/",
									HttpMethod:   rshttp.MethodGet,
									ContentType:  rshttp.MIMEApplicationJSON,
									Timeout:      0,
									Assertion: models.AssertionV1{
										StatusCode: http.StatusOK,
									},
									Created:      time.Now(),
									LastModified: time.Now(),
								},
							},
							Created:      time.Now(),
							LastModified: time.Now(),
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
			manager := &webServiceScheduleManager{
				webServiceSchedulers: tt.fields.webServiceSchedulers,
				resultChan:           tt.fields.resultChan,
				closeChan:            make(chan bool, 1),
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
