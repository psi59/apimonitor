package repositories

import (
	"github.com/realsangil/apimonitor/models"
	"github.com/realsangil/apimonitor/pkg/rsdb"
	"github.com/realsangil/apimonitor/pkg/rsmodel"
)

type WebServiceTestResultRepository interface {
	rsdb.Repository
	GetResultList(conn rsdb.Connection, webServiceTest *models.WebServiceTest, request models.WebServiceTestResultListRequest) (*rsmodel.PaginatedList, error)
}

type WebServiceTestResultRepositoryImp struct {
	rsdb.Repository
}

func (repository *WebServiceTestResultRepositoryImp) CreateTable(conn rsdb.Connection) error {
	m := &models.WebServiceTestResult{}
	if err := conn.Conn().AutoMigrate(m).Error; err != nil {
		return rsdb.HandleSQLError(err)
	}
	if err := conn.Conn().Model(m).
		AddForeignKey("web_service_test_id", "web_service_tests(id)", "CASCADE", "CASCADE").Error; err != nil {
		return rsdb.HandleSQLError(err)
	}
	if err := conn.Conn().Model(m).AddIndex("idx_tested_at_status_code_is_success", "tested_at", "status_code", "is_success").Error; err != nil {
		return rsdb.HandleSQLError(err)
	}
	return nil
}

func (repository *WebServiceTestResultRepositoryImp) GetResultList(
	conn rsdb.Connection, webServiceTest *models.WebServiceTest, request models.WebServiceTestResultListRequest) (*rsmodel.PaginatedList, error) {
	sql := conn.Conn().Table("web_service_test_results AS wstr")
	sql = sql.Joins("INNER JOIN web_service_tests AS wst ON wstr.web_service_test_id=wst.id AND wst.id=?", webServiceTest.Id)

	query := rsdb.NewEmptyQuery()
	if !request.IsSuccess.IsBoth() {
		q, _ := rsdb.NewQuery("wstr.is_success = ?", request.IsSuccess)
		query = query.And(q)
	}

	switch {
	case !request.StartTestedAt.IsZero():
		q, _ := rsdb.NewQuery("wstr.tested_at > ?", request.StartTestedAt)
		query = query.And(q)
	case !request.EndTestedAt.IsZero():
		q, _ := rsdb.NewQuery("wstr.tested_at < ?", request.EndTestedAt)
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

	items := make([]*models.WebServiceTestResult, 0)
	if err := sql.Order("wstr.tested_at DESC").Offset(listFilter.Offset()).Limit(listFilter.NumItem).Find(&items).Error; err != nil {
		return nil, rsdb.HandleSQLError(err)
	}

	return &rsmodel.PaginatedList{
		CurrentPage: request.Page,
		NumItem:     request.NumItem,
		TotalCount:  totalCount,
		Items:       items,
	}, nil
}

func NewWebServiceTestResultRepository() WebServiceTestResultRepository {
	return &WebServiceTestResultRepositoryImp{
		&rsdb.DefaultRepository{},
	}
}
