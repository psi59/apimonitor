package services

import (
	"github.com/realsangil/apimonitor/models"
	"github.com/realsangil/apimonitor/pkg/amerr"
	"github.com/realsangil/apimonitor/pkg/rsdb"
	"github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/rslog"
	"github.com/realsangil/apimonitor/pkg/rsvalid"
	"github.com/realsangil/apimonitor/repositories"
)

type WebServiceService interface {
	CreateWebService(transaction rsdb.Transaction, request models.WebServiceRequest) (*models.WebService, *amerr.ErrorWithLanguage)
	GetWebServiceById(transaction rsdb.Transaction, webService *models.WebService) *amerr.ErrorWithLanguage
	DeleteWebServiceById(transaction rsdb.Transaction, webService *models.WebService) *amerr.ErrorWithLanguage
	UpdateWebServiceById(transaction rsdb.Transaction, webService *models.WebService, request models.WebServiceRequest) *amerr.ErrorWithLanguage
}

type WebServiceServiceImpl struct {
	webServiceRepository repositories.WebServiceRepository
}

func (service *WebServiceServiceImpl) CreateWebService(transaction rsdb.Transaction, request models.WebServiceRequest) (*models.WebService, *amerr.ErrorWithLanguage) {
	if rsvalid.IsZero(transaction, request) {
		return nil, amerr.GetErrorsFromCode(amerr.ErrInternalServer)
	}

	webService, err := models.NewWebService(request)
	if err != nil {
		rslog.Error(err)
		return nil, amerr.GetErrorsFromCode(amerr.ErrBadRequest)
	}

	if err := service.webServiceRepository.Create(transaction, webService); err != nil {
		rslog.Error(err)
		switch err {
		case rsdb.ErrInvalidData:
			return nil, amerr.GetErrorsFromCode(amerr.ErrBadRequest)
		case rsdb.ErrDuplicateData:
			return nil, amerr.GetErrorsFromCode(amerr.ErrDuplicatedWebService)
		}
		return nil, amerr.GetErrorsFromCode(amerr.ErrInternalServer)
	}

	return webService, nil
}

func (service *WebServiceServiceImpl) GetWebServiceById(transaction rsdb.Transaction, webService *models.WebService) *amerr.ErrorWithLanguage {
	if rsvalid.IsZero(transaction, webService) {
		return amerr.GetErrorsFromCode(amerr.ErrInternalServer)
	}

	if err := service.webServiceRepository.GetById(transaction, webService); err != nil {
		switch err {
		case rsdb.ErrRecordNotFound:
			return amerr.GetErrorsFromCode(amerr.ErrWebServiceNotFound)
		}
		rslog.Error(err)
		return amerr.GetErrInternalServer()
	}

	return nil
}

func (service *WebServiceServiceImpl) DeleteWebServiceById(transaction rsdb.Transaction, webService *models.WebService) *amerr.ErrorWithLanguage {
	if rsvalid.IsZero(transaction, webService) {
		return amerr.GetErrorsFromCode(amerr.ErrInternalServer)
	}

	if err := service.GetWebServiceById(transaction, webService); err != nil {
		return err
	}

	if err := service.webServiceRepository.DeleteById(transaction, webService); err != nil {
		switch err {
		case rsdb.ErrRecordNotFound:
			return amerr.GetErrorsFromCode(amerr.ErrWebServiceNotFound)
		}
		rslog.Error(err)
		return amerr.GetErrInternalServer()
	}

	return nil
}

func (service *WebServiceServiceImpl) UpdateWebServiceById(transaction rsdb.Transaction, webService *models.WebService, request models.WebServiceRequest) *amerr.ErrorWithLanguage {
	if rsvalid.IsZero(transaction, webService) {
		return amerr.GetErrorsFromCode(amerr.ErrInternalServer)
	}

	if err := service.GetWebServiceById(transaction, webService); err != nil {
		return err
	}

	if err := webService.UpdateFromRequest(request); err != nil {
		rslog.Error(err)
		return amerr.GetErrorsFromCode(amerr.ErrBadRequest)
	}

	if err := service.webServiceRepository.Save(transaction, webService); err != nil {
		switch err {
		case rsdb.ErrRecordNotFound:
			return amerr.GetErrorsFromCode(amerr.ErrWebServiceNotFound)
		}
		rslog.Error(err)
		return amerr.GetErrInternalServer()
	}

	return nil
}

func NewWebServiceService(webServiceRepository repositories.WebServiceRepository) (WebServiceService, error) {
	if rsvalid.IsZero(webServiceRepository) {
		return nil, rserrors.ErrInvalidParameter
	}
	return &WebServiceServiceImpl{
		webServiceRepository: webServiceRepository,
	}, nil
}
