package services

import (
	"github.com/pkg/errors"
	"github.com/realsangil/apimonitor/models"
	"github.com/realsangil/apimonitor/pkg/amerr"
	"github.com/realsangil/apimonitor/pkg/rsdb"
	"github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/rslog"
	"github.com/realsangil/apimonitor/pkg/rsvalid"
	"github.com/realsangil/apimonitor/repositories"
)

type EndpointService interface {
	CreateEndpoint(webService *models.WebService, request models.EndpointRequest) (*models.Endpoint, *amerr.ErrorWithLanguage)
}

type EndpointServiceImpl struct {
	endpointRepository repositories.EndpointRepository
}

func (service *EndpointServiceImpl) CreateEndpoint(webService *models.WebService, request models.EndpointRequest) (*models.Endpoint, *amerr.ErrorWithLanguage) {
	if rsvalid.IsZero(webService, request) {
		rslog.Error(rserrors.ErrInvalidParameter)
		return nil, amerr.GetErrInternalServer()
	}

	endpoint, err := models.NewEndpoint(webService, request)
	if err != nil {
		rslog.Error(err)
		return nil, amerr.GetErrorsFromCode(amerr.ErrBadRequest)
	}

	if err := service.endpointRepository.Create(rsdb.GetConnection(), endpoint); err != nil {
		switch err {
		case rsdb.ErrForeignKeyConstraint:
			return nil, amerr.GetErrorsFromCode(amerr.ErrWebServiceNotFound)
		case rsdb.ErrDuplicateData:
			return nil, amerr.GetErrorsFromCode(amerr.ErrDuplicatedEndpoint)
		default:
			rslog.Error(err)
			return nil, amerr.GetErrInternalServer()
		}
	}

	return endpoint, nil
}

func NewEndpointService(endpointRepository repositories.EndpointRepository) (EndpointService, error) {
	if rsvalid.IsZero(endpointRepository) {
		return nil, errors.Wrap(rserrors.ErrInvalidParameter, "EndpointService")
	}
	return &EndpointServiceImpl{
		endpointRepository: endpointRepository,
	}, nil
}