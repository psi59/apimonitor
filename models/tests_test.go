package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/realsangil/apimonitor/pkg/rshttp"
	"github.com/realsangil/apimonitor/pkg/rsjson"
	"github.com/realsangil/apimonitor/pkg/rsmodels"
	"github.com/realsangil/apimonitor/pkg/testutils"
)

func TestTest_Validate(t *testing.T) {
	type fields struct {
		DefaultValidateChecker rsmodels.DefaultValidateChecker
		Id                     int64
		WebServiceId           int64
		Path                   rshttp.EndpointPath
		HttpMethod             rshttp.Method
		ContentType            rshttp.ContentType
		RequestData            rsjson.MapJson
		Header                 rshttp.Header
		QueryParam             rshttp.Query
		Created                time.Time
		LastModified           time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "pass",
			fields: fields{
				WebServiceId: 1,
				Path:         "/v1/test",
				HttpMethod:   "GET",
				ContentType:  "application/json",
				Created:      time.Now(),
				LastModified: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "invalid path",
			fields: fields{
				WebServiceId: 1,
				Path:         "/?/test",
				HttpMethod:   "GET",
				ContentType:  "application/json",
				Created:      time.Now(),
				LastModified: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "invalid path",
			fields: fields{
				WebServiceId: 1,
				// Path:         "/v1/test",
				HttpMethod:   "GET",
				ContentType:  "application/json",
				Created:      time.Now(),
				LastModified: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "invalid parameter",
			fields: fields{
				WebServiceId: 1,
				Path:         "/v1/test",
				HttpMethod:   "invalid",
				ContentType:  "application/json",
				Created:      time.Now(),
				LastModified: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "invalid parameter",
			fields: fields{
				WebServiceId: 1,
				Path:         "/v1/test",
				HttpMethod:   "GET",
				ContentType:  "invalid",
				Created:      time.Now(),
				LastModified: time.Now(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := &Test{
				DefaultValidateChecker: tt.fields.DefaultValidateChecker,
				Id:                     tt.fields.Id,
				WebServiceId:           tt.fields.WebServiceId,
				Path:                   tt.fields.Path,
				HttpMethod:             tt.fields.HttpMethod,
				ContentType:            tt.fields.ContentType,
				RequestData:            tt.fields.RequestData,
				Header:                 tt.fields.Header,
				QueryParam:             tt.fields.QueryParam,
				Created:                tt.fields.Created,
				LastModified:           tt.fields.LastModified,
			}
			if err := test.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTest_UpdateFromRequest(t *testing.T) {
	testutils.MonkeyAll()

	mockRequestData := rsjson.MapJson{
		"key1": "value1",
		"key2": 2,
	}

	mockHeader := rshttp.Header{
		"Authorization":   "Bearer access_token",
		"accept-language": "ko",
	}

	type args struct {
		request TestRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				request: TestRequest{
					Path:        "/path/to/file",
					HttpMethod:  rshttp.MethodGet,
					ContentType: rshttp.MIMEApplicationJSON,
					RequestData: mockRequestData,
					Header:      mockHeader,
					QueryParam:  nil,
				},
			},
			wantErr: false,
		},
		{
			name: "pass",
			args: args{
				request: TestRequest{
					Path:        "/path/to/file",
					HttpMethod:  rshttp.MethodGet,
					ContentType: rshttp.MIMEApplicationJSON,
					RequestData: mockRequestData,
					Header:      mockHeader,
					QueryParam:  nil,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid request",
			args: args{
				request: TestRequest{
					Path:        "/???/asdas",
					HttpMethod:  rshttp.MethodGet,
					ContentType: rshttp.MIMEApplicationJSON,
					RequestData: mockRequestData,
					Header:      mockHeader,
					QueryParam:  nil,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid request",
			args: args{
				request: TestRequest{
					Path: "/path/to/file",
					// HttpMethod:  http.MethodGet,
					ContentType: rshttp.MIMEApplicationJSON,
					RequestData: mockRequestData,
					Header:      mockHeader,
					QueryParam:  nil,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid request",
			args: args{
				request: TestRequest{
					Path:       "/path/to/file",
					HttpMethod: rshttp.MethodGet,
					// ContentType: http.MIMEApplicationJSON,
					RequestData: mockRequestData,
					Header:      mockHeader,
					QueryParam:  nil,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := &Test{
				WebServiceId: 1,
				Created:      time.Now(),
			}
			if err := test.UpdateFromRequest(tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("UpdateFromRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewTest(t *testing.T) {
	testutils.MonkeyAll()

	webService := &WebService{
		Id: 1,
	}

	request := TestRequest{
		Path:        "/path/to/file",
		HttpMethod:  rshttp.MethodGet,
		ContentType: rshttp.MIMEApplicationJSON,
	}

	type args struct {
		webService *WebService
		request    TestRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *Test
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				webService: webService,
				request:    request,
			},
			want: &Test{
				DefaultValidateChecker: rsmodels.ValidatedDefaultValidateChecker,
				WebServiceId:           1,
				Path:                   request.Path,
				HttpMethod:             request.HttpMethod,
				ContentType:            request.ContentType,
				RequestData:            request.RequestData,
				Header:                 request.Header,
				QueryParam:             request.QueryParam,
				Created:                time.Now(),
				LastModified:           time.Now(),
			},
			wantErr: false,
		},
		{
			name: "invalid args",
			args: args{
				webService: webService,
				// request:    request,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid args",
			args: args{
				// webService: webService,
				request: request,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid request",
			args: args{
				webService: webService,
				request: TestRequest{
					Path: "//??/asd",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTest(tt.args.webService, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEndpoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTestRequest_Validate(t *testing.T) {
	type fields struct {
		Path        rshttp.EndpointPath
		HttpMethod  rshttp.Method
		ContentType rshttp.ContentType
		RequestData rsjson.MapJson
		Header      rshttp.Header
		QueryParam  rshttp.Query
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "pass",
			fields: fields{
				Path:        "/path/to/file",
				HttpMethod:  rshttp.MethodGet,
				ContentType: rshttp.MIMEApplicationJSON,
				RequestData: nil,
				Header:      nil,
				QueryParam:  nil,
			},
			wantErr: false,
		},
		{
			name: "invalid parameter",
			fields: fields{
				// Path:        "/path/to/file",
				HttpMethod:  rshttp.MethodGet,
				ContentType: rshttp.MIMEApplicationJSON,
				RequestData: nil,
				Header:      nil,
				QueryParam:  nil,
			},
			wantErr: true,
		},
		{
			name: "invalid parameter",
			fields: fields{
				Path:        "/???/to/file",
				HttpMethod:  rshttp.MethodGet,
				ContentType: rshttp.MIMEApplicationJSON,
				RequestData: nil,
				Header:      nil,
				QueryParam:  nil,
			},
			wantErr: true,
		},
		{
			name: "invalid parameter",
			fields: fields{
				Path:        "/path/to/file",
				HttpMethod:  "invalid",
				ContentType: rshttp.MIMEApplicationJSON,
				RequestData: nil,
				Header:      nil,
				QueryParam:  nil,
			},
			wantErr: true,
		},
		{
			name: "invalid parameter",
			fields: fields{
				Path:        "/path/to/file",
				HttpMethod:  rshttp.MethodGet,
				ContentType: "invalid",
				RequestData: nil,
				Header:      nil,
				QueryParam:  nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := TestRequest{
				Path:        tt.fields.Path,
				HttpMethod:  tt.fields.HttpMethod,
				ContentType: tt.fields.ContentType,
				RequestData: tt.fields.RequestData,
				Header:      tt.fields.Header,
				QueryParam:  tt.fields.QueryParam,
			}
			if err := e.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
