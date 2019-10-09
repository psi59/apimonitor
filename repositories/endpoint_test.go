package repositories

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/realsangil/apimonitor/models"
	"github.com/realsangil/apimonitor/pkg/rsdb"
	"github.com/realsangil/apimonitor/pkg/rshttp"
	"github.com/realsangil/apimonitor/pkg/testutils"
)

var countColumn = []string{
	"count",
}

var endpointListItemColumns = []string{
	"id", "web_service_id", "path", "http_method", "desc", "created", "last_modified",
}

func TestEndpointRepositoryImpl_GetList(t *testing.T) {
	testutils.MonkeyAll()

	countRows := sqlmock.NewRows(countColumn).AddRow(int64(1))
	rows := sqlmock.NewRows(endpointListItemColumns).AddRow(int64(1), int64(1), "/test", rshttp.MethodGet, "", time.Now(), time.Now())

	type args struct {
		items  *[]*models.EndpointListItem
		filter rsdb.ListFilter
		orders rsdb.Orders
	}
	tests := []struct {
		name     string
		args     args
		mockFunc rsdb.MockFunc
		want     int
		wantErr  error
	}{
		{
			name: "pass",
			args: args{
				items: &[]*models.EndpointListItem{},
				filter: rsdb.ListFilter{
					Page:    1,
					NumItem: 20,
					Conditions: map[string]interface{}{
						"web_service_id": int64(1),
					},
				},
				orders: rsdb.Orders{},
			},
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT count\(\*\)`).WithArgs(int64(1)).WillReturnRows(countRows)
				mock.ExpectQuery("SELECT").WithArgs(int64(1)).WillReturnRows(rows)
				mock.ExpectQuery(`SELECT \* FROM "web_services"`).WithArgs(int64(1)).WillReturnRows(rows)
			},
			want:    1,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock, err := rsdb.CreateMockDB()
			if err != nil {
				t.Fatal(err)
			}
			tt.mockFunc(mock)

			repository := &EndpointRepositoryImpl{
				Repository: &rsdb.DefaultRepository{},
			}
			conn := rsdb.NewConnection(gormDB)

			got, err := repository.GetList(conn, tt.args.items, tt.args.filter, tt.args.orders)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
