package repositories

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/realsangil/apimonitor/models"
	"github.com/realsangil/apimonitor/pkg/rsdb"
	"github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/rshttp"
	"github.com/realsangil/apimonitor/pkg/testutils"
)

var (
	webServiceColumns = []string{"id", "host", "http_schema", "desc", "favicon", "schedule", "created", "last_modified"}
	testsColums       = []string{
		"id",
		"web_service_id",
		"path",
		"http_method",
		"content_type",
		"desc",
		"request_data",
		"header",
		"query_param",
		"timeout",
		"assertion",
		"created",
		"last_modified",
	}
)

func TestWebServiceRepositoryImpl_GetAllWebServicesWithTests(t *testing.T) {
	testutils.MonkeyAll()

	getWebServiceRows := func() *sqlmock.Rows {
		webServiceRows := sqlmock.NewRows(webServiceColumns)
		webServiceRows.
			AddRow(1, "github.com", "https", "깃허브", "", "1m", time.Now(), time.Now()).
			AddRow(2, "google.com", "https", "구글", "", "5m", time.Now(), time.Now())
		return webServiceRows
	}

	getGithubTestRows := func() *sqlmock.Rows {
		githubTestRows := sqlmock.NewRows(testsColums)
		githubTestRows.AddRow(1, 1, "/", rshttp.MethodGet, rshttp.MIMEApplicationJSON, "", nil, nil, nil, 0, nil, time.Now(), time.Now())
		return githubTestRows
	}

	githubTests := []models.Test{
		{
			Id:           1,
			WebServiceId: 1,
			Path:         "/",
			Method:       rshttp.MethodGet,
			ContentType:  rshttp.MIMEApplicationJSON,
			Description:  "",
			RequestData:  nil,
			Header:       nil,
			QueryParam:   nil,
			Timeout:      0,
			Assertion:    models.AssertionV1{},
			CreatedAt:    time.Now(),
			ModifiedAt:   time.Now(),
		},
	}

	webServices := []models.WebService{
		{
			Id:          1,
			Host:        "github.com",
			Schema:      "https",
			Description: "깃허브",
			Favicon:     "",
			Schedule:    models.ScheduleOneMinute,
			Tests:       githubTests,
			CreatedAt:   time.Now(),
			ModifiedAt:  time.Now(),
		},
		{
			Id:          2,
			Host:        "google.com",
			Schema:      "https",
			Description: "구글",
			Favicon:     "",
			Schedule:    models.ScheduleFiveMinute,
			Tests:       []models.Test{},
			CreatedAt:   time.Now(),
			ModifiedAt:  time.Now(),
		},
	}

	tests := []struct {
		name     string
		mockFunc rsdb.MockFunc
		want     []models.WebService
		wantErr  error
	}{
		{
			name: "pass",
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM "web_services"`).WillReturnRows(getWebServiceRows())
				mock.ExpectQuery(`SELECT (.+) FROM "tests"`).WillReturnRows(getGithubTestRows())
			},
			want:    webServices,
			wantErr: nil,
		},
		{
			name: "failed_web_services_query",
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM "web_services"`).WillReturnError(rserrors.ErrUnexpected)
			},
			want:    nil,
			wantErr: rserrors.ErrUnexpected,
		},
		{
			name: "failed_web_services_tests_query",
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "web_services"`).WillReturnRows(getWebServiceRows())
				mock.ExpectQuery(`SELECT(.+)"tests"`).WillReturnError(rserrors.ErrUnexpected)
			},
			want:    nil,
			wantErr: rserrors.ErrUnexpected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock, err := rsdb.CreateMockDB()
			if err != nil {
				t.Fatal(err)
			}
			repository := &WebServiceRepositoryImpl{Repository: &rsdb.DefaultRepository{}}
			tt.mockFunc(mock)
			got, err := repository.GetAllWebServicesWithTests(rsdb.NewConnection(gormDB))
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
