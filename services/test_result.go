package services

import (
	"github.com/realsangil/apimonitor/models"
	"github.com/realsangil/apimonitor/pkg/amerr"
	"github.com/realsangil/apimonitor/pkg/rsdb"
	"github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/rsmodels"
	"github.com/realsangil/apimonitor/pkg/rsvalid"
	"github.com/realsangil/apimonitor/repositories"
)

type TestResultService interface {
	GetResultListByTestId(test *models.Test, request models.TestResultListRequest) (*rsmodels.PaginatedList, *amerr.ErrorWithLanguage)
}

type TestResultServiceImpl struct {
	testResultRepository repositories.TestResultRepository
}

func (service *TestResultServiceImpl) GetResultListByTestId(test *models.Test, request models.TestResultListRequest) (*rsmodels.PaginatedList, *amerr.ErrorWithLanguage) {
	list, err := service.testResultRepository.GetResultList(rsdb.GetConnection(), test, request)
	if err != nil {
		switch err {
		case rsdb.ErrForeignKeyConstraint:
			return nil, amerr.GetErrorsFromCode(amerr.ErrTestNotFound)
		default:
			return nil, amerr.GetErrInternalServer()
		}
	}
	return list, nil
}

func NewTestResultService(testResultRepository repositories.TestResultRepository) (TestResultService, error) {
	if rsvalid.IsZero(testResultRepository) {
		return nil, rserrors.ErrInvalidParameter
	}
	return &TestResultServiceImpl{
		testResultRepository: testResultRepository,
	}, nil
}
