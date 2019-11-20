package repositories

import (
	"github.com/realsangil/apimonitor/models"
	"github.com/realsangil/apimonitor/pkg/rsdb"
	"github.com/realsangil/apimonitor/pkg/rsmodels"
)

type TestResultRepository interface {
	rsdb.Repository
	GetResultListByTest(conn rsdb.Connection, test *models.Test, request models.TestResultListRequest) (*rsmodels.PaginatedList, error)
	GetResultListByWebService(conn rsdb.Connection, webService *models.WebService, request models.TestResultListRequest) (*rsmodels.PaginatedList, error)
}

type TestResultRepositoryImp struct {
	rsdb.Repository
}

func (repository *TestResultRepositoryImp) CreateTable(conn rsdb.Connection) error {
	m := &models.TestResult{}
	if err := conn.Conn().AutoMigrate(m).Error; err != nil {
		return rsdb.HandleSQLError(err)
	}
	if err := conn.Conn().Model(m).
		AddForeignKey("test_id", "tests(id)", "CASCADE", "CASCADE").Error; err != nil {
		return rsdb.HandleSQLError(err)
	}
	if err := conn.Conn().Model(m).AddIndex("idx_tested_at_status_code_is_success", "tested_at", "status_code", "is_success").Error; err != nil {
		return rsdb.HandleSQLError(err)
	}
	return nil
}

func (repository *TestResultRepositoryImp) GetResultListByWebService(
	conn rsdb.Connection, webService *models.WebService, request models.TestResultListRequest) (*rsmodels.PaginatedList, error) {
	return repository.getResultList(conn, webService, request)
}

func (repository *TestResultRepositoryImp) GetResultListByTest(
	conn rsdb.Connection, test *models.Test, request models.TestResultListRequest) (*rsmodels.PaginatedList, error) {
	return repository.getResultList(conn, test, request)
}

func (repository *TestResultRepositoryImp) getResultList(conn rsdb.Connection, joinObject interface{}, request models.TestResultListRequest) (*rsmodels.PaginatedList, error) {
	sql := conn.Conn().Table("test_results AS tr").Select("*")

	switch obj := joinObject.(type) {
	case *models.WebService:
		sql = sql.Joins("INNER JOIN tests AS t ON tr.test_id=t.id AND t.web_service_id=?", obj.Id)
	case *models.Test:
		sql = sql.Joins("INNER JOIN tests AS t ON tr.test_id=t.id AND t.id=?", obj.Id)
	}

	query := rsdb.NewEmptyQuery()
	if !request.IsSuccess.IsBoth() {
		q, _ := rsdb.NewQuery("tr.is_success = ?", request.IsSuccess)
		query = query.And(q)
	}

	switch {
	case !request.StartTestedAt.IsZero():
		q, _ := rsdb.NewQuery("tr.tested_at > ?", request.StartTestedAt)
		query = query.And(q)
	case !request.EndTestedAt.IsZero():
		q, _ := rsdb.NewQuery("tr.tested_at < ?", request.EndTestedAt)
		query = query.And(q)
	}

	sql = sql.Where(query.Where(), query.Values()...)

	listFilter := rsdb.ListFilter{
		Page:       request.Page,
		NumItem:    request.NumItem,
		Conditions: nil,
	}

	var totalCount int
	if err := sql.Count(&totalCount).Error; err != nil {
		return nil, rsdb.HandleSQLError(err)
	}

	items := make([]*models.TestResult, 0)
	if err := sql.Order("tr.tested_at DESC").Offset(listFilter.Offset()).Limit(listFilter.NumItem).Find(&items).Error; err != nil {
		return nil, rsdb.HandleSQLError(err)
	}

	return &rsmodels.PaginatedList{
		CurrentPage: request.Page,
		NumItem:     request.NumItem,
		TotalCount:  totalCount,
		Items:       items,
	}, nil

}

func NewTestResultRepository() TestResultRepository {
	return &TestResultRepositoryImp{
		&rsdb.DefaultRepository{},
	}
}
