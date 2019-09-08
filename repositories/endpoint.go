package repositories

import (
	"github.com/pkg/errors"

	"github.com/realsangil/apimonitor/models"
	"github.com/realsangil/apimonitor/pkg/rsdb"
)

type EndpointRepository interface {
	rsdb.Repository
	GetByIdAndWebServiceId(conn rsdb.Connection, endpoint *models.Endpoint) error
}

type EndpointRepositoryImpl struct {
	rsdb.Repository
}

func (repository *EndpointRepositoryImpl) GetByIdAndWebServiceId(conn rsdb.Connection, endpoint *models.Endpoint) error {
	if err := conn.Conn().
		Where("web_service_id=? AND id=?", endpoint.WebServiceId, endpoint.Id).
		First(endpoint).Error; err != nil {
		return rsdb.HandleSQLError(err)
	}

	return nil
}

func (repository EndpointRepositoryImpl) CreateTable(transaction rsdb.Connection) error {
	m := &models.Endpoint{}
	tx := transaction.Conn()
	if tx.HasTable(m) {
		return nil
	}
	if err := tx.AutoMigrate(m).Error; err != nil {
		return errors.WithStack(err)
	}
	if err := tx.Model(m).AddForeignKey("web_service_id", "web_services(id)", "CASCADE", "CASCADE").Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func NewEndpointRepository() EndpointRepository {
	return &EndpointRepositoryImpl{&rsdb.DefaultRepository{}}
}
