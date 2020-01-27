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

var _ TestService = &TestServiceImpl{}

type TestService interface {
	CreateTest(webService *models.WebService, request models.TestRequest) (*models.Test, *amerr.ErrorWithLanguage)
	GetTestById(endpoint *models.Test) *amerr.ErrorWithLanguage
	DeleteTestById(endpoint *models.Test) *amerr.ErrorWithLanguage
	GetTestList(request models.TestListRequest) (*rsmodels.PaginatedList, *amerr.ErrorWithLanguage)
	UpdateTestById(test *models.Test, request models.TestRequest) *amerr.ErrorWithLanguage
	ExecuteTest(test *models.Test) *amerr.ErrorWithLanguage
}

type TestServiceImpl struct {
	testRepository      repositories.TestRepository
	testScheduleManager ScheduleManager
}

func (service *TestServiceImpl) CreateTest(webService *models.WebService, request models.TestRequest) (*models.Test, *amerr.ErrorWithLanguage) {
	if rsvalid.IsZero(webService, request) {
		rslog.Error(rserrors.ErrInvalidParameter)
		return nil, amerr.GetErrInternalServer()
	}

	test, err := models.NewTest(webService, request)
	if err != nil {
		rslog.Error(err)
		return nil, amerr.GetErrorsFromCode(amerr.ErrBadRequest)
	}

	if err := service.testRepository.Create(rsdb.GetConnection(), test); err != nil {
		switch err {
		case rsdb.ErrForeignKeyConstraint:
			return nil, amerr.GetErrorsFromCode(amerr.ErrWebServiceNotFound)
		case rsdb.ErrDuplicateData:
			return nil, amerr.GetErrorsFromCode(amerr.ErrDuplicatedTest)
		default:
			rslog.Error(err)
			return nil, amerr.GetErrInternalServer()
		}
	}
	test.WebService = webService

	if err := service.testScheduleManager.AddSchedule(test); err != nil {
		rslog.Error(err)
		return nil, amerr.GetErrInternalServer()
	}

	return test, nil
}

func (service *TestServiceImpl) GetTestById(test *models.Test) *amerr.ErrorWithLanguage {
	if rsvalid.IsZero(test) {
		rslog.Error(errors.Wrap(rserrors.ErrInvalidParameter, "Test"))
		return amerr.GetErrInternalServer()
	}

	if rsvalid.IsZero(test.Id) {
		rslog.Error(errors.Wrap(rserrors.ErrInvalidParameter, "Test.Id"))
		return amerr.GetErrInternalServer()
	}

	if err := service.testRepository.GetById(rsdb.GetConnection(), test); err != nil {
		switch err {
		case rsdb.ErrRecordNotFound:
			return amerr.GetErrorsFromCode(amerr.ErrTestNotFound)
		default:
			rslog.Error(err)
			return amerr.GetErrInternalServer()
		}
	}

	return nil
}

func (service *TestServiceImpl) DeleteTestById(test *models.Test) *amerr.ErrorWithLanguage {
	if rsvalid.IsZero(test) {
		rslog.Error(errors.Wrap(rserrors.ErrInvalidParameter, "Test"))
		return amerr.GetErrInternalServer()
	}

	if rsvalid.IsZero(test.Id) {
		rslog.Error(errors.Wrap(rserrors.ErrInvalidParameter, "Test.Id"))
		return amerr.GetErrInternalServer()
	}

	if err := service.testRepository.DeleteById(rsdb.GetConnection(), test); err != nil {
		rslog.Error(err)
		return amerr.GetErrInternalServer()
	}

	if err := service.testScheduleManager.RemoveSchedule(test); err != nil {
		rslog.Error(err)
		return amerr.GetErrInternalServer()
	}

	return nil
}

func (service *TestServiceImpl) GetTestList(request models.TestListRequest) (*rsmodels.PaginatedList, *amerr.ErrorWithLanguage) {
	filter := rsdb.ListFilter{
		Page:    request.Page,
		NumItem: request.NumItem,
		Conditions: map[string]interface{}{
			"web_service_id": request.WebServiceId,
		},
	}

	items := make([]*models.TestListItem, 0)
	totalCount, err := service.testRepository.GetList(rsdb.GetConnection(), &items, filter, rsdb.Orders{
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

func (service *TestServiceImpl) UpdateTestById(test *models.Test, request models.TestRequest) *amerr.ErrorWithLanguage {
	if err := service.testRepository.GetById(rsdb.GetConnection(), test); err != nil {
		switch err {
		case rsdb.ErrRecordNotFound:
			return amerr.GetErrorsFromCode(amerr.ErrTestNotFound)
		default:
			rslog.Error(err)
			return amerr.GetErrInternalServer()
		}
	}
	_ = test.UpdateFromRequest(request)

	if err := service.testRepository.Save(rsdb.GetConnection(), test); err != nil {
		switch err {
		case rsdb.ErrRecordNotFound:
			return amerr.GetErrorsFromCode(amerr.ErrTestNotFound)
		case rsdb.ErrInvalidData:
			return amerr.GetErrorsFromCode(amerr.ErrBadRequest)
		default:
			rslog.Error(err)
			return amerr.GetErrInternalServer()
		}
	}

	if err := service.testScheduleManager.UpdateSchedule(test); err != nil {
		rslog.Error(err)
		return amerr.GetErrInternalServer()
	}

	return nil
}

func (service *TestServiceImpl) ExecuteTest(test *models.Test) *amerr.ErrorWithLanguage {
	if err := service.testRepository.GetById(rsdb.GetConnection(), test); err != nil {
		switch err {
		case rsdb.ErrRecordNotFound:
			return amerr.GetErrorsFromCode(amerr.ErrTestNotFound)
		default:
			rslog.Error(err)
			return amerr.GetErrInternalServer()
		}
	}
	service.testScheduleManager.ExecuteSchedule(test)
	return nil
}

func NewTestService(testRepository repositories.TestRepository, testScheduleManager ScheduleManager) (TestService, error) {
	if rsvalid.IsZero(testRepository, testScheduleManager) {
		return nil, errors.Wrap(rserrors.ErrInvalidParameter, "TestService")
	}
	return &TestServiceImpl{
		testRepository:      testRepository,
		testScheduleManager: testScheduleManager,
	}, nil
}
