package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/realsangil/apimonitor/pkg/http"
	"github.com/realsangil/apimonitor/pkg/rsjson"
	"github.com/realsangil/apimonitor/pkg/testutils"
)

func TestEndpoint_Validate(t *testing.T) {
	type fields struct {
		DefaultValidateChecker DefaultValidateChecker
		Id                     int64
		WebServiceId           int64
		Path                   http.EndpointPath
		HttpMethod             http.Method
		ContentType            http.ContentType
		RequestData            rsjson.MapJson
		Header                 rsjson.MapJson
		QueryParam             rsjson.MapJson
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			endpoint := &Endpoint{
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
			if err := endpoint.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEndpoint_UpdateFromRequest(t *testing.T) {
	testutils.MonkeyAll()

	mockRequestData := rsjson.MapJson{
		"key1": "value1",
		"key2": 2,
	}

	mockHeader := rsjson.MapJson{
		"Authorization":   "Bearer access_token",
		"accept-language": "ko",
	}

	type args struct {
		request EndpointRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				request: EndpointRequest{
					Path:        "/path/to/file",
					HttpMethod:  http.MethodGet,
					ContentType: http.MIMEApplicationJSON,
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
				request: EndpointRequest{
					Path:        "/path/to/file",
					HttpMethod:  http.MethodGet,
					ContentType: http.MIMEApplicationJSON,
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
				request: EndpointRequest{
					Path:        "/???/asdas",
					HttpMethod:  http.MethodGet,
					ContentType: http.MIMEApplicationJSON,
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
				request: EndpointRequest{
					Path: "/path/to/file",
					// HttpMethod:  http.MethodGet,
					ContentType: http.MIMEApplicationJSON,
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
				request: EndpointRequest{
					Path:       "/path/to/file",
					HttpMethod: http.MethodGet,
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
			endpoint := &Endpoint{
				WebServiceId: 1,
				Created:      time.Now(),
			}
			if err := endpoint.UpdateFromRequest(tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("UpdateFromRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewEndpoint(t *testing.T) {
	testutils.MonkeyAll()

	webService := &WebService{
		Id: 1,
	}

	request := EndpointRequest{
		Path:        "/path/to/file",
		HttpMethod:  http.MethodGet,
		ContentType: http.MIMEApplicationJSON,
	}

	type args struct {
		webService *WebService
		request    EndpointRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *Endpoint
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				webService: webService,
				request:    request,
			},
			want: &Endpoint{
				DefaultValidateChecker: ValidatedDefaultValidateChecker,
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
				request: EndpointRequest{
					Path: "//??/asd",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEndpoint(tt.args.webService, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEndpoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}