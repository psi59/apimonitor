package models

import (
	"reflect"
	"testing"
	"time"
)

func Test_hostRegexpFindStringSubmatch(t *testing.T) {
	type args struct {
		host string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				host: "realsangil.github.io",
			},
			want: []string{
				"realsangil.github.io",
				"",
				"realsangil.github.io",
				"realsangil.github.io",
				"",
			},
			wantErr: false,
		},
		{
			name: "pass",
			args: args{
				host: "https://realsangil.github.io",
			},
			want: []string{
				"https://realsangil.github.io",
				"https",
				"realsangil.github.io",
				"realsangil.github.io",
				"",
			},
			wantErr: false,
		},
		{
			name: "pass",
			args: args{
				host: "http://realsangil.github.io",
			},
			want: []string{
				"http://realsangil.github.io",
				"http",
				"realsangil.github.io",
				"realsangil.github.io",
				"",
			},
			wantErr: false,
		},
		{
			name: "invalid host",
			args: args{
				host: "asdasdasd",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid host",
			args: args{
				host: "asdasdasdzxczxczxc.a",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hostRegexpFindStringSubmatch(tt.args.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("hostRegexpFindStringSubmatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("hostRegexpFindStringSubmatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebService_UpdateFromRequest(t *testing.T) {
	type args struct {
		request WebServiceRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				request: WebServiceRequest{
					Host:    "https://realsangil.github.io",
					Favicon: "",
					Desc:    "sangil's blog",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid host",
			args: args{
				request: WebServiceRequest{
					Host:    "",
					Favicon: "",
					Desc:    "sangil's blog",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			webService := &WebService{}
			if err := webService.UpdateFromRequest(tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("WebService.UpdateFromRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWebService_Validate(t *testing.T) {
	type fields struct {
		DefaultValidateChecker DefaultValidateChecker
		Id                     int64
		Host                   string
		HttpSchema             string
		Desc                   string
		Favicon                string
		Created                time.Time
		LastModified           time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			webService := &WebService{
				DefaultValidateChecker: tt.fields.DefaultValidateChecker,
				Id:                     tt.fields.Id,
				Host:                   tt.fields.Host,
				HttpSchema:             tt.fields.HttpSchema,
				Desc:                   tt.fields.Desc,
				Favicon:                tt.fields.Favicon,
				Created:                tt.fields.Created,
				LastModified:           tt.fields.LastModified,
			}
			if err := webService.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("WebService.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewWebService(t *testing.T) {
	type args struct {
		request WebServiceRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *WebService
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewWebService(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWebService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWebService() = %v, want %v", got, tt.want)
			}
		})
	}
}
