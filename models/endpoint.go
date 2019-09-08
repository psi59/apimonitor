package models

import (
	"time"

	"github.com/pkg/errors"

	"github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/rshttp"
	"github.com/realsangil/apimonitor/pkg/rsjson"
	"github.com/realsangil/apimonitor/pkg/rsvalid"
)

type Endpoint struct {
	DefaultValidateChecker
	Id           int64               `json:"id"`
	WebServiceId int64               `json:"-"`
	Path         rshttp.EndpointPath `json:"path"`
	HttpMethod   rshttp.Method       `json:"http_method"`
	ContentType  rshttp.ContentType  `json:"content_type"`
	RequestData  rsjson.MapJson      `json:"request_data" gorm:"Type:JSON"`
	Header       rsjson.MapJson      `json:"header" gorm:"Type:JSON"`
	QueryParam   rsjson.MapJson      `json:"query_param" gorm:"Type:JSON"`
	Created      time.Time           `json:"created"`
	LastModified time.Time           `json:"last_modified"`
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
	if err := endpoint.HttpMethod.Validate(); err != nil {
		return errors.WithStack(err)
	}
	if err := endpoint.ContentType.Validate(); err != nil {
		return err
	}
	endpoint.SetValidated()
	return nil
}

type EndpointRequest struct {
	Path        rshttp.EndpointPath `json:"path"`
	HttpMethod  rshttp.Method       `json:"http_method"`
	ContentType rshttp.ContentType  `json:"content_type"`
	RequestData rsjson.MapJson      `json:"request_data" gorm:"Type:JSON"`
	Header      rsjson.MapJson      `json:"header" gorm:"Type:JSON"`
	QueryParam  rsjson.MapJson      `json:"query_param" gorm:"Type:JSON"`
}

func (e EndpointRequest) Validate() error {
	if rsvalid.IsZero(
		e.Path,
		e.HttpMethod,
		e.ContentType,
	) {
		return errors.Wrap(rserrors.ErrInvalidParameter, "EndpointRequest")
	}
	if err := e.Path.Validate(); err != nil {
		return errors.WithStack(err)
	}
	if err := e.HttpMethod.Validate(); err != nil {
		return errors.WithStack(err)
	}
	if err := e.ContentType.Validate(); err != nil {
		return err
	}
	return nil
}
