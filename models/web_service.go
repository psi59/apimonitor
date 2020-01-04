package models

import (
	"regexp"
	"time"

	"github.com/pkg/errors"

	"github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/rsmodels"
	"github.com/realsangil/apimonitor/pkg/rsstr"
	"github.com/realsangil/apimonitor/pkg/rsvalid"
)

var regexExtractHost = regexp.MustCompile(`^(?:(?:(https?)?(?:\:?\/\/))|(?:\/\/))?(((?:\w{1,100}\.)?\w{2,300}\.\w{2,100})(\.\w{2,100})*)`)

type WebService struct {
	rsmodels.DefaultValidateChecker
	Id          string    `json:"id" gorm:"private_key"`
	Host        string    `json:"host" gorm:"unique"`
	Schema      string    `json:"schema" gorm:"Size:20;Default:'http'"`
	Description string    `json:"description" gorm:"Type:TEXT"`
	CreatedAt   time.Time `json:"createdAt"`
	ModifiedAt  time.Time `json:"modifiedAt"`
}

func (webService *WebService) Validate() error {
	if rsvalid.IsZero(
		webService.Id, webService.ModifiedAt, webService.CreatedAt, webService.Host,
	) {
		return errors.Wrap(rserrors.ErrInvalidParameter, "webService")
	}
	webService.SetValidated()
	return nil
}

func (webService *WebService) UpdateFromRequest(request WebServiceRequest) error {
	host, err := hostRegexpFindStringSubmatch(request.Host)
	if err != nil {
		return errors.WithStack(err)
	}

	webService.Host = host[2]
	webService.Schema = host[1]
	webService.Description = request.Description
	webService.ModifiedAt = time.Now()

	return nil
}

func NewWebService(request WebServiceRequest) (*WebService, error) {
	webService := &WebService{
		CreatedAt: time.Now(),
	}

	webService.Id = rsstr.NewUUID()
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
	Host        string `json:"host"`
	Description string `json:"description"`
}

type WebServiceListRequest struct {
	Page          int
	NumItem       int
	SearchKeyword string
}
