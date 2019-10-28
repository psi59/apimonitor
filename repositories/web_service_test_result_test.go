package repositories

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/realsangil/apimonitor/models"
	"github.com/realsangil/apimonitor/pkg/rsdb"
	"github.com/realsangil/apimonitor/pkg/rsmodel"
	"github.com/realsangil/apimonitor/pkg/testutils"
)

var webServiceTestResultColumn = []string{
	"id", "web_service_test_id", "is_success", "status_code", "response_time", "tested_at",
}

func TestWebServiceTestResultRepositoryImp_GetResultList(t *testing.T) {
	testutils.MonkeyAll()

	webServiceTest := &models.WebServiceTest{Id: 1}

	type args struct {
		webServiceTest *models.WebServiceTest
		request        models.WebServiceTestResultListRequest
	}
	tests := []struct {
		name     string
		args     args
		mockFunc rsdb.MockFunc
		want     *rsmodel.PaginatedList
		wantErr  error
	}{
		{
			name: "pass",
			args: args{
				webServiceTest: webServiceTest,
				request: models.WebServiceTestResultListRequest{
					Page:      1,
					NumItem:   1,
					IsSuccess: "",
					// StartTestedAt: time.Time{},
					// EndTestedAt:   time.Time{},
				},
			},
			mockFunc: func(mock sqlmock.Sqlmock) {
				countRows := sqlmock.NewRows(countColumn).AddRow(1)
				mock.ExpectQuery(`SELECT count\(\*\)`).WithArgs(1).WillReturnRows(countRows)
				rows := sqlmock.NewRows(webServiceTestResultColumn).
					AddRow("web_service_test_result_0", 1, 1, 200, 1, time.Now())
				mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(rows)
			},
			want: &rsmodel.PaginatedList{
				CurrentPage: 1,
				NumItem:     1,
				TotalCount:  1,
				Items: []*models.WebServiceTestResult{
					{
						Id:               "web_service_test_result_0",
						WebServiceTestId: 1,
						IsSuccess:        true,
						StatusCode:       200,
						ResponseTime:     1,
						TestedAt:         time.Now(),
					},
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &WebServiceTestResultRepositoryImp{rsdb.NewDefaultRepository()}
			gormDB, mock, err := rsdb.CreateMockDB()
			if err != nil {
				t.Fatal(err)
			}
			conn := rsdb.NewConnection(gormDB)
			tt.mockFunc(mock)
			got, err := repository.GetResultList(conn, tt.args.webServiceTest, tt.args.request)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
