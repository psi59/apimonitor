package services

import (
	"github.com/realsangil/apimonitor/pkg/amerr"
	"github.com/realsangil/apimonitor/pkg/rsmodel"
	"github.com/realsangil/apimonitor/repositories"
)

type WebServiceTestResultService interface {
	GetResultList() (*rsmodel.PaginatedList, *amerr.ErrorWithLanguage)
}

type WebServiceTestResultServiceImpl struct {
	webServiceTestResultRepository repositories.WebServiceTestResultRepository
}

func (service *WebServiceTestResultServiceImpl) GetResultList() (*rsmodel.PaginatedList, *amerr.ErrorWithLanguage) {
	panic("implement me")
}
