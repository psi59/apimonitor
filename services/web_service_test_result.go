package services

import (
	"github.com/realsangil/apimonitor/models"
	"github.com/realsangil/apimonitor/pkg/amerr"
	"github.com/realsangil/apimonitor/pkg/rsdb"
	"github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/rsmodel"
	"github.com/realsangil/apimonitor/pkg/rsvalid"
	"github.com/realsangil/apimonitor/repositories"
)

type WebServiceTestResultService interface {
	GetResultListByTestId(webServiceTest *models.WebServiceTest, request models.WebServiceTestResultListRequest) (*rsmodel.PaginatedList, *amerr.ErrorWithLanguage)
}

type WebServiceTestResultServiceImpl struct {
	webServiceTestResultRepository repositories.WebServiceTestResultRepository
}

func (service *WebServiceTestResultServiceImpl) GetResultListByTestId(webServiceTest *models.WebServiceTest, request models.WebServiceTestResultListRequest) (*rsmodel.PaginatedList, *amerr.ErrorWithLanguage) {
	list, err := service.webServiceTestResultRepository.GetResultList(rsdb.GetConnection(), webServiceTest, request)
	if err != nil {
		switch err {
		case rsdb.ErrForeignKeyConstraint:
			return nil, amerr.GetErrorsFromCode(amerr.ErrWebServiceTestNotFound)
		default:
			return nil, amerr.GetErrInternalServer()
		}
	}
	return list, nil
}

func NewWebServiceTestResultService(webServiceTestResultRepository repositories.WebServiceTestResultRepository) (WebServiceTestResultService, error) {
	if rsvalid.IsZero(webServiceTestResultRepository) {
		return nil, rserrors.ErrInvalidParameter
	}
	return &WebServiceTestResultServiceImpl{
		webServiceTestResultRepository: webServiceTestResultRepository,
	}, nil
}
