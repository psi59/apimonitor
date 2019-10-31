package services

import (
	"github.com/pkg/errors"
	"github.com/realsangil/apimonitor/models"
	"github.com/realsangil/apimonitor/pkg/amerr"
	"github.com/realsangil/apimonitor/pkg/rsdb"
	"github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/rslog"
	"github.com/realsangil/apimonitor/pkg/rsmodels"
	"github.com/realsangil/apimonitor/pkg/rsvalid"
	"github.com/realsangil/apimonitor/repositories"
)

type WebServiceTestService interface {
	CreateWebServiceTest(webService *models.WebService, request models.WebServiceTestRequest) (*models.WebServiceTest, *amerr.ErrorWithLanguage)
	GetWebServiceTestById(endpoint *models.WebServiceTest) *amerr.ErrorWithLanguage
	DeleteWebServiceTestById(endpoint *models.WebServiceTest) *amerr.ErrorWithLanguage
	GetWebServiceTestList(request models.WebServiceTestListRequest) (*rsmodels.PaginatedList, *amerr.ErrorWithLanguage)
}

type WebServiceTestServiceImpl struct {
	webServiceTestRepository repositories.WebServiceTestRepository
}

func (service *WebServiceTestServiceImpl) CreateWebServiceTest(webService *models.WebService, request models.WebServiceTestRequest) (*models.WebServiceTest, *amerr.ErrorWithLanguage) {
	if rsvalid.IsZero(webService, request) {
		rslog.Error(rserrors.ErrInvalidParameter)
		return nil, amerr.GetErrInternalServer()
	}

	webServiceTest, err := models.NewWebServiceTest(webService, request)
	if err != nil {
		rslog.Error(err)
		return nil, amerr.GetErrorsFromCode(amerr.ErrBadRequest)
	}

	if err := service.webServiceTestRepository.Create(rsdb.GetConnection(), webServiceTest); err != nil {
		switch err {
		case rsdb.ErrForeignKeyConstraint:
			return nil, amerr.GetErrorsFromCode(amerr.ErrWebServiceNotFound)
		case rsdb.ErrDuplicateData:
			return nil, amerr.GetErrorsFromCode(amerr.ErrDuplicatedWebServiceTest)
		default:
			rslog.Error(err)
			return nil, amerr.GetErrInternalServer()
		}
	}

	return webServiceTest, nil
}

func (service *WebServiceTestServiceImpl) GetWebServiceTestById(webServiceTest *models.WebServiceTest) *amerr.ErrorWithLanguage {
	if rsvalid.IsZero(webServiceTest) {
		rslog.Error(errors.Wrap(rserrors.ErrInvalidParameter, "WebServiceTest"))
		return amerr.GetErrInternalServer()
	}

	if rsvalid.IsZero(webServiceTest.Id) {
		rslog.Error(errors.Wrap(rserrors.ErrInvalidParameter, "WebServiceTest.Id"))
		return amerr.GetErrInternalServer()
	}

	if err := service.webServiceTestRepository.GetById(rsdb.GetConnection(), webServiceTest); err != nil {
		switch err {
		case rsdb.ErrRecordNotFound:
			return amerr.GetErrorsFromCode(amerr.ErrWebServiceTestNotFound)
		default:
			rslog.Error(err)
			return amerr.GetErrInternalServer()
		}
	}

	return nil
}

func (service *WebServiceTestServiceImpl) DeleteWebServiceTestById(webServiceTest *models.WebServiceTest) *amerr.ErrorWithLanguage {
	if rsvalid.IsZero(webServiceTest) {
		rslog.Error(errors.Wrap(rserrors.ErrInvalidParameter, "WebServiceTest"))
		return amerr.GetErrInternalServer()
	}

	if rsvalid.IsZero(webServiceTest.Id) {
		rslog.Error(errors.Wrap(rserrors.ErrInvalidParameter, "WebServiceTest.Id"))
		return amerr.GetErrInternalServer()
	}

	if err := service.webServiceTestRepository.DeleteById(rsdb.GetConnection(), webServiceTest); err != nil {
		rslog.Error(err)
		return amerr.GetErrInternalServer()
	}

	return nil
}

func (service *WebServiceTestServiceImpl) GetWebServiceTestList(request models.WebServiceTestListRequest) (*rsmodels.PaginatedList, *amerr.ErrorWithLanguage) {
	filter := rsdb.ListFilter{
		Page:    request.Page,
		NumItem: request.NumItem,
		Conditions: map[string]interface{}{
			"web_service_id": request.WebServiceId,
		},
	}

	items := make([]*models.WebServiceTestListItem, 0)
	totalCount, err := service.webServiceTestRepository.GetList(rsdb.GetConnection(), &items, filter, rsdb.Orders{
		rsdb.Order{
			Field: "path",
			IsASC: true,
		},
	})

	if err != nil {
		return nil, amerr.GetErrInternalServer()
	}

	return &rsmodels.PaginatedList{
		CurrentPage: request.Page,
		NumItem:     request.NumItem,
		TotalCount:  totalCount,
		Items:       items,
	}, nil
}

func NewWebServiceTestService(webServiceTestRepository repositories.WebServiceTestRepository) (WebServiceTestService, error) {
	if rsvalid.IsZero(webServiceTestRepository) {
		return nil, errors.Wrap(rserrors.ErrInvalidParameter, "WebServiceTestService")
	}
	return &WebServiceTestServiceImpl{
		webServiceTestRepository: webServiceTestRepository,
	}, nil
}
