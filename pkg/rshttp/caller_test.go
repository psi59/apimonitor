package rshttp

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	type args struct {
		request Request
	}
	tests := []struct {
		name    string
		args    args
		want    Response
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				request: Request{
					RawUrl: "https://dev.ccbblab.com/api/v1/posts",
					Header: nil,
					Query:  nil,
					Body:   nil,
				},
			},
			want: &HttpResponse{
				StatusCode:   200,
				ResponseTime: 0,
				Body:         nil,
			},
			wantErr: false,
		},
		{
			name: "pass",
			args: args{
				request: Request{
					RawUrl: "https://dev.ccbblab.com/api/v1/postss",
					Header: nil,
					Query:  nil,
					Body:   nil,
				},
			},
			want: &HttpResponse{
				StatusCode:   http.StatusNotFound,
				ResponseTime: 0,
				Body:         nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want.GetStatusCode(), got.GetStatusCode())
		})
	}
}
