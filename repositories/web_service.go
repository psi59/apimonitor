package repositories

import (
	"github.com/pkg/errors"

	"github.com/realsangil/apimonitor/models"
	"github.com/realsangil/apimonitor/pkg/rsdb"
)

type WebServiceRepository interface {
	rsdb.Repository
}

type WebServiceRepositoryImpl struct {
	rsdb.Repository
}

func (repository WebServiceRepositoryImpl) CreateTable(transaction rsdb.Connection) error {
	m := &models.WebService{}
	tx := transaction.Conn()
	if tx.HasTable(m) {
		return nil
	}
	if err := tx.AutoMigrate(m).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func NewWebServiceRepository() WebServiceRepository {
	return &WebServiceRepositoryImpl{&rsdb.DefaultRepository{}}
}
