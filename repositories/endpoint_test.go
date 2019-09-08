package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/realsangil/apimonitor/models"
	"github.com/realsangil/apimonitor/pkg/rsdb"
	"github.com/realsangil/apimonitor/pkg/testutils"
)

func TestEndpointRepositoryImpl_GetList(t *testing.T) {
	testutils.MonkeyAll()

	type args struct {
		items  *[]*models.EndpointListItem
		filter rsdb.ListFilter
		order  rsdb.Order
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
				order: rsdb.Order{},
			},
			mockFunc: nil,
			want:     0,
			wantErr:  nil,
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

			got, err := repository.GetList(conn, tt.args.items, tt.args.filter, tt.args.order)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
