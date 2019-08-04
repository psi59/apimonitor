package models

import (
	"regexp"
	"time"

	"github.com/pkg/errors"

	"github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/rsvalid"
)

var regexExtractHost = regexp.MustCompile(`^(?:(?:(https?)?(?:\:?\/\/))|(?:\/\/))?(((?:\w{1,100}\.)?\w{2,300}\.\w{2,100})(\.\w{2,100})*)`)

type WebService struct {
	Id           int64     `json:"id" gorm:"private_key"`
	Host         string    `json:"host" gorm:"unique"`
	HttpSchema   string    `json:"http_schema" gorm:"Size:20;Default:'http'"`
	Desc         string    `json:"desc" gorm:"Type:TEXT"`
	Favicon      string    `json:"favicon" gorm:"Type:Text"`
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"last_modified"`
	isValidated  bool      `gorm:"-"`
}

func (webService *WebService) IsValidated() bool {
	return webService.isValidated
}

func (webService *WebService) Validate() error {
	if rsvalid.IsZero(
		webService.LastModified, webService.Created,
		webService.Desc, webService.Host,
	) {
		return errors.Wrap(rserrors.ErrInvalidParameter, "webService")
	}
	webService.isValidated = true
	return nil
}

func NewWebService(request WebServiceRequest) (*WebService, error) {
	if !regexExtractHost.MatchString(request.Host) {
		return nil, rserrors.ErrInvalidParameter
	}
	host := regexExtractHost.FindStringSubmatch(request.Host)
	webService := &WebService{
		Host:         host[2],
		HttpSchema:   host[1],
		Desc:         request.Desc,
		Favicon:      request.Favicon,
		Created:      time.Now(),
		LastModified: time.Now(),
	}

	if err := webService.Validate(); err != nil {
		return nil, errors.WithStack(err)
	}

	return webService, nil
}

type WebServiceRequest struct {
	Host    string `json:"host"`
	Desc    string `json:"desc"`
	Favicon string `json:"favicon"`
}
