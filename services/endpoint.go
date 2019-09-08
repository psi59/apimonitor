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
	GetEndpointById(endpoint *models.Endpoint) *amerr.ErrorWithLanguage
	DeleteEndpointById(endpoint *models.Endpoint) *amerr.ErrorWithLanguage
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

func (service *EndpointServiceImpl) GetEndpointById(endpoint *models.Endpoint) *amerr.ErrorWithLanguage {
	if rsvalid.IsZero(endpoint) {
		rslog.Error(errors.Wrap(rserrors.ErrInvalidParameter, "Endpoint"))
		return amerr.GetErrInternalServer()
	}

	if rsvalid.IsZero(endpoint.Id) {
		rslog.Error(errors.Wrap(rserrors.ErrInvalidParameter, "Endpoint.Id"))
		return amerr.GetErrInternalServer()
	}

	if err := service.endpointRepository.GetById(rsdb.GetConnection(), endpoint); err != nil {
		switch err {
		case rsdb.ErrRecordNotFound:
			return amerr.GetErrorsFromCode(amerr.ErrEndpointNotFound)
		default:
			rslog.Error(err)
			return amerr.GetErrInternalServer()
		}
	}

	return nil
}

func (service *EndpointServiceImpl) DeleteEndpointById(endpoint *models.Endpoint) *amerr.ErrorWithLanguage {
	if rsvalid.IsZero(endpoint) {
		rslog.Error(errors.Wrap(rserrors.ErrInvalidParameter, "Endpoint"))
		return amerr.GetErrInternalServer()
	}

	if rsvalid.IsZero(endpoint.Id) {
		rslog.Error(errors.Wrap(rserrors.ErrInvalidParameter, "Endpoint.Id"))
		return amerr.GetErrInternalServer()
	}

	if err := service.endpointRepository.DeleteById(rsdb.GetConnection(), endpoint); err != nil {
		rslog.Error(err)
		return amerr.GetErrInternalServer()
	}

	return nil
}

func NewEndpointService(endpointRepository repositories.EndpointRepository) (EndpointService, error) {
	if rsvalid.IsZero(endpointRepository) {
		return nil, errors.Wrap(rserrors.ErrInvalidParameter, "EndpointService")
	}
	return &EndpointServiceImpl{
		endpointRepository: endpointRepository,
	}, nil
}
