package rshttp

import (
	"testing"
)

func TestEndpointPath_String(t *testing.T) {
	tests := []struct {
		name         string
		endpointPath EndpointPath
		want         string
	}{
		{
			name:         "pass",
			endpointPath: "/test",
			want:         "/test",
		},
		{
			name:         "pass",
			endpointPath: "/test/path",
			want:         "/test/path",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.endpointPath.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndpointPath_Validate(t *testing.T) {
	tests := []struct {
		name         string
		endpointPath EndpointPath
		wantErr      bool
	}{
		{
			name:         "pass",
			endpointPath: "/test",
			wantErr:      false,
		},
		{
			name:         "invalid path",
			endpointPath: "??/?",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.endpointPath.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHttpMethod_Validate(t *testing.T) {
	tests := []struct {
		name    string
		method  Method
		wantErr bool
	}{
		{
			name:    "get",
			method:  MethodGet,
			wantErr: false,
		},
		{
			name:    "head",
			method:  MethodHead,
			wantErr: false,
		},
		{
			name:    "post",
			method:  MethodPost,
			wantErr: false,
		},
		{
			name:    "put",
			method:  MethodPut,
			wantErr: false,
		},
		{
			name:    "patch",
			method:  MethodPatch,
			wantErr: false,
		},
		{
			name:    "delete",
			method:  MethodDelete,
			wantErr: false,
		},
		{
			name:    "connect",
			method:  MethodConnect,
			wantErr: false,
		},
		{
			name:    "options",
			method:  MethodOptions,
			wantErr: false,
		},
		{
			name:    "trace",
			method:  MethodTrace,
			wantErr: false,
		},
		{
			name:    "invalid",
			method:  "invalid",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.method.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMethod_String(t *testing.T) {
	tests := []struct {
		name   string
		method Method
		want   string
	}{
		{
			name:   "pass",
			method: "GET",
			want:   "GET",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.method.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContentType_String(t *testing.T) {
	tests := []struct {
		name        string
		contentType ContentType
		want        string
	}{
		{
			name:        "pass",
			contentType: "text/json",
			want:        "text/json",
		},
		{
			name:        "pass",
			contentType: "application/xml",
			want:        "application/xml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.contentType.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContentType_Validate(t *testing.T) {
	tests := []struct {
		name        string
		contentType ContentType
		wantErr     bool
	}{
		{
			name:        "pass",
			contentType: "application/xml",
			wantErr:     false,
		},
		{
			name:        "pass",
			contentType: "application/xml; utf-8",
			wantErr:     false,
		},
		{
			name:        "pass",
			contentType: "text/javascript",
			wantErr:     false,
		},
		{
			name:        "pass",
			contentType: "multipart/form-data",
			wantErr:     false,
		},
		{
			name:        "invalid contentType",
			contentType: "application",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.contentType.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
