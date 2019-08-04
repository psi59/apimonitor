package models

import (
	"time"

	"github.com/realsangil/apimonitor/pkg/rsjson"
)

type Endpoint struct {
	Id           int64          `json:"id"`
	WebServiceId int64          `json:"-"`
	Path         string         `json:"path"`
	HttpMethod   string         `json:"http_method"`
	ContentType  string         `json:"content_type"`
	RequestData  rsjson.MapJson `json:"request_data" gorm:"Type:JSON"`
	Header       rsjson.MapJson `json:"header" gorm:"Type:JSON"`
	QueryParam   rsjson.MapJson `json:"query_param" gorm:"Type:JSON"`
	Created      time.Time      `json:"created"`
	LastModified time.Time      `json:"last_modified"`
}
