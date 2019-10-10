package repositories

import (
	"github.com/pkg/errors"

	"github.com/realsangil/apimonitor/models"
	"github.com/realsangil/apimonitor/pkg/rsdb"
	"github.com/realsangil/apimonitor/pkg/rsvalid"
)

type WebServiceTestRepository interface {
	rsdb.Repository
	GetByIdAndWebServiceId(conn rsdb.Connection, endpoint *models.WebServiceTest) error
	GetList(conn rsdb.Connection, items *[]*models.WebServiceTestListItem, filter rsdb.ListFilter, orders rsdb.Orders) (int, error)
}

type WebServiceTestRepositoryImpl struct {
	rsdb.Repository
}

func (repository *WebServiceTestRepositoryImpl) GetByIdAndWebServiceId(conn rsdb.Connection, endpoint *models.WebServiceTest) error {
	if err := conn.Conn().
		Where("web_service_id=? AND id=?", endpoint.WebServiceId, endpoint.Id).
		First(endpoint).Error; err != nil {
		return rsdb.HandleSQLError(err)
	}

	return nil
}

func (repository *WebServiceTestRepositoryImpl) GetList(conn rsdb.Connection, items *[]*models.WebServiceTestListItem, filter rsdb.ListFilter, orders rsdb.Orders) (int, error) {
	where := rsdb.NewEmptyQuery()
	if v, exist := filter.Conditions["web_service_id"]; exist {
		w, _ := rsdb.NewQuery("web_service_id=?", v)
		where.And(w)
	}

	query := conn.Conn().Model(items).Where(where.Where(), where.Values()...)
	if !rsvalid.IsZero(orders) {
		query = query.Order(orders.String())
	}

	var totalCount int
	if err := query.Count(&totalCount).Error; err != nil {
		return 0, rsdb.HandleSQLError(err)
	}

	if err := query.Preload("WebService").Offset(filter.Offset()).Limit(filter.NumItem).Find(items).Error; err != nil {
		return 0, rsdb.HandleSQLError(err)
	}

	return totalCount, nil
}

func (repository WebServiceTestRepositoryImpl) CreateTable(transaction rsdb.Connection) error {
	m := &models.WebServiceTest{}
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

func NewWebServiceTestRepository() WebServiceTestRepository {
	return &WebServiceTestRepositoryImpl{&rsdb.DefaultRepository{}}
}
