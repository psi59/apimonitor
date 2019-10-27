package repositories

import (
	"github.com/realsangil/apimonitor/models"
	"github.com/realsangil/apimonitor/pkg/rsdb"
)

type WebServiceTestResultRepository interface {
	rsdb.Repository
}

type WebServiceTestResultRepositoryImp struct {
	rsdb.Repository
}

func (repository *WebServiceTestResultRepositoryImp) CreateTable(conn rsdb.Connection) error {
	m := &models.WebServiceTestResult{}
	if err := conn.Conn().AutoMigrate(m).Error; err != nil {
		return rsdb.HandleSQLError(err)
	}
	if err := conn.Conn().
		AddForeignKey("web_service_test_id", "web_service_tests(id)", "CASCADE", "CASCADE").Error; err != nil {
		return rsdb.HandleSQLError(err)
	}
	return nil
}

func NewWebServiceTestResultRepository() WebServiceTestResultRepository {
	return &WebServiceTestResultRepositoryImp{
		&rsdb.DefaultRepository{},
	}
}
