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
	DefaultValidateChecker
	Id           int64     `json:"id" gorm:"private_key"`
	Host         string    `json:"host" gorm:"unique"`
	HttpSchema   string    `json:"http_schema" gorm:"Size:20;Default:'http'"`
	Desc         string    `json:"desc" gorm:"Type:TEXT"`
	Favicon      string    `json:"favicon" gorm:"Type:Text"`
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"last_modified"`
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

func (webService *WebService) UpdateFromRequest(request WebServiceRequest) error {
	host, err := hostRegexpFindStringSubmatch(request.Host)
	if err != nil {
		return errors.WithStack(err)
	}

	webService.Host = host[2]
	webService.HttpSchema = host[1]
	webService.Desc = request.Desc
	webService.Favicon = request.Favicon
	webService.LastModified = time.Now()

	return nil
}

func NewWebService(request WebServiceRequest) (*WebService, error) {
	webService := &WebService{
		Created: time.Now(),
	}

	if err := webService.UpdateFromRequest(request); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := webService.Validate(); err != nil {
		return nil, errors.WithStack(err)
	}

	return webService, nil
}

func hostRegexpFindStringSubmatch(host string) ([]string, error) {
	if !regexExtractHost.MatchString(host) {
		return nil, rserrors.ErrInvalidParameter
	}
	return regexExtractHost.FindStringSubmatch(host), nil
}

type WebServiceRequest struct {
	Host    string `json:"host"`
	Desc    string `json:"desc"`
	Favicon string `json:"favicon"`
}
