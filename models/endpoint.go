package models

import (
	"time"

	"github.com/pkg/errors"

	"github.com/realsangil/apimonitor/pkg/http"
	"github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/rsjson"
	"github.com/realsangil/apimonitor/pkg/rsvalid"
)

type Endpoint struct {
	DefaultValidateChecker
	Id           int64             `json:"id"`
	WebServiceId int64             `json:"-"`
	Path         http.EndpointPath `json:"path"`
	HttpMethod   http.Method       `json:"http_method"`
	ContentType  http.ContentType  `json:"content_type"`
	RequestData  rsjson.MapJson    `json:"request_data" gorm:"Type:JSON"`
	Header       rsjson.MapJson    `json:"header" gorm:"Type:JSON"`
	QueryParam   rsjson.MapJson    `json:"query_param" gorm:"Type:JSON"`
	Created      time.Time         `json:"created"`
	LastModified time.Time         `json:"last_modified"`
}

func NewEndpoint(webService *WebService, request EndpointRequest) (*Endpoint, error) {
	if rsvalid.IsZero(webService, request) {
		return nil, errors.Wrap(rserrors.ErrInvalidParameter, "Endpoint")
	}
	endpoint := &Endpoint{
		WebServiceId: webService.Id,
		Created:      time.Now(),
	}
	if err := endpoint.UpdateFromRequest(request); err != nil {
		return nil, errors.WithStack(err)
	}
	return endpoint, nil
}

func (endpoint *Endpoint) UpdateFromRequest(request EndpointRequest) error {
	endpoint.Path = request.Path
	endpoint.HttpMethod = request.HttpMethod
	endpoint.ContentType = request.ContentType
	endpoint.RequestData = request.RequestData
	endpoint.Header = request.Header
	endpoint.QueryParam = request.Header
	endpoint.LastModified = time.Now()
	return endpoint.Validate()
}

func (endpoint *Endpoint) Validate() error {
	if rsvalid.IsZero(
		endpoint.WebServiceId,
		endpoint.Path,
		endpoint.HttpMethod,
		endpoint.ContentType,
		endpoint.Created,
		endpoint.LastModified,
	) {
		return errors.Wrap(rserrors.ErrInvalidParameter, "Endpoint")
	}
	if err := endpoint.Path.Validate(); err != nil {
		return errors.WithStack(err)
	}
	endpoint.SetValidated()
	return nil
}

type EndpointRequest struct {
	Path        http.EndpointPath `json:"path"`
	HttpMethod  http.Method       `json:"http_method"`
	ContentType http.ContentType  `json:"content_type"`
	RequestData rsjson.MapJson    `json:"request_data" gorm:"Type:JSON"`
	Header      rsjson.MapJson    `json:"header" gorm:"Type:JSON"`
	QueryParam  rsjson.MapJson    `json:"query_param" gorm:"Type:JSON"`
}
