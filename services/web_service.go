package services

import (
	"github.com/realsangil/apimonitor/models"
	"github.com/realsangil/apimonitor/pkg/amerr"
	"github.com/realsangil/apimonitor/pkg/rsdb"
	"github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/rslog"
	"github.com/realsangil/apimonitor/pkg/rsmodel"
	"github.com/realsangil/apimonitor/pkg/rsvalid"
	"github.com/realsangil/apimonitor/repositories"
)

type WebServiceService interface {
	CreateWebService(request models.WebServiceRequest) (*models.WebService, *amerr.ErrorWithLanguage)
	GetWebServiceById(webService *models.WebService) *amerr.ErrorWithLanguage
	DeleteWebServiceById(webService *models.WebService) *amerr.ErrorWithLanguage
	UpdateWebServiceById(webService *models.WebService, request models.WebServiceRequest) *amerr.ErrorWithLanguage
	GetWebServiceList(request models.WebServiceListRequest) (*rsmodel.PaginatedList, *amerr.ErrorWithLanguage)
}

type WebServiceServiceImpl struct {
	webServiceRepository repositories.WebServiceRepository
}

func (service *WebServiceServiceImpl) CreateWebService(request models.WebServiceRequest) (*models.WebService, *amerr.ErrorWithLanguage) {
	if rsvalid.IsZero(request) {
		return nil, amerr.GetErrorsFromCode(amerr.ErrInternalServer)
	}

	webService, err := models.NewWebService(request)
	if err != nil {
		rslog.Error(err)
		return nil, amerr.GetErrorsFromCode(amerr.ErrBadRequest)
	}

	if err := service.webServiceRepository.Create(rsdb.GetConnection(), webService); err != nil {
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

func (service *WebServiceServiceImpl) GetWebServiceById(webService *models.WebService) *amerr.ErrorWithLanguage {
	if rsvalid.IsZero(webService) {
		return amerr.GetErrorsFromCode(amerr.ErrInternalServer)
	}

	if err := service.webServiceRepository.GetById(rsdb.GetConnection(), webService); err != nil {
		switch err {
		case rsdb.ErrRecordNotFound:
			return amerr.GetErrorsFromCode(amerr.ErrWebServiceNotFound)
		}
		rslog.Error(err)
		return amerr.GetErrInternalServer()
	}

	return nil
}

func (service *WebServiceServiceImpl) DeleteWebServiceById(webService *models.WebService) *amerr.ErrorWithLanguage {
	if rsvalid.IsZero(webService) {
		return amerr.GetErrorsFromCode(amerr.ErrInternalServer)
	}

	if err := service.GetWebServiceById(webService); err != nil {
		return err
	}

	if err := service.webServiceRepository.DeleteById(rsdb.GetConnection(), webService); err != nil {
		switch err {
		case rsdb.ErrRecordNotFound:
			return amerr.GetErrorsFromCode(amerr.ErrWebServiceNotFound)
		}
		rslog.Error(err)
		return amerr.GetErrInternalServer()
	}

	return nil
}

func (service *WebServiceServiceImpl) UpdateWebServiceById(webService *models.WebService, request models.WebServiceRequest) *amerr.ErrorWithLanguage {
	if rsvalid.IsZero(webService) {
		return amerr.GetErrorsFromCode(amerr.ErrInternalServer)
	}

	if err := service.GetWebServiceById(webService); err != nil {
		return err
	}

	if err := webService.UpdateFromRequest(request); err != nil {
		rslog.Error(err)
		return amerr.GetErrorsFromCode(amerr.ErrBadRequest)
	}

	if err := service.webServiceRepository.Save(rsdb.GetConnection(), webService); err != nil {
		switch err {
		case rsdb.ErrRecordNotFound:
			return amerr.GetErrorsFromCode(amerr.ErrWebServiceNotFound)
		}
		rslog.Error(err)
		return amerr.GetErrInternalServer()
	}

	return nil
}

func (service *WebServiceServiceImpl) GetWebServiceList(request models.WebServiceListRequest) (*rsmodel.PaginatedList, *amerr.ErrorWithLanguage) {
	if rsvalid.IsZero(request) {
		return nil, amerr.GetErrorsFromCode(amerr.ErrInternalServer)
	}

	items := make([]*models.WebService, 0)
	totalCount, err := service.webServiceRepository.List(rsdb.GetConnection(), &items, rsdb.ListFilter{
		Page:       request.Page,
		NumItem:    request.NumItem,
		Conditions: map[string]interface{}{},
	}, rsdb.Orders{
		{
			Field: "host",
			IsASC: true,
		},
	})

	if err != nil {
		rslog.Error(err)
		return nil, amerr.GetErrInternalServer()
	}

	return &rsmodel.PaginatedList{
		TotalCount:  totalCount,
		CurrentPage: request.Page,
		NumItem:     request.NumItem,
		Items:       items,
	}, nil
}

func NewWebServiceService(webServiceRepository repositories.WebServiceRepository) (WebServiceService, error) {
	if rsvalid.IsZero(webServiceRepository) {
		return nil, rserrors.ErrInvalidParameter
	}
	return &WebServiceServiceImpl{
		webServiceRepository: webServiceRepository,
	}, nil
}
