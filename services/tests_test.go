package services

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/realsangil/apimonitor/models"
	"github.com/realsangil/apimonitor/pkg/amerr"
	"github.com/realsangil/apimonitor/pkg/rsdb"
	mocks2 "github.com/realsangil/apimonitor/pkg/rsdb/mocks"
	"github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/rshttp"
	"github.com/realsangil/apimonitor/pkg/rsjson"
	"github.com/realsangil/apimonitor/pkg/rsmodels"
	"github.com/realsangil/apimonitor/pkg/testutils"
	"github.com/realsangil/apimonitor/repositories"
	"github.com/realsangil/apimonitor/repositories/mocks"
)

type testMockFunc func(mockConn *mocks2.Connection, mockTestRepository *mocks.TestRepository)

func TestNewTestService(t *testing.T) {
	type args struct {
		testRepository repositories.TestRepository
	}
	tests := []struct {
		name    string
		args    args
		want    TestService
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				testRepository: &mocks.TestRepository{},
			},
			want: &TestServiceImpl{
				testRepository: &mocks.TestRepository{},
			},
			wantErr: false,
		},
		{
			name: "invalid parameter",
			args: args{
				// testRepository: &mocks.TestRepository{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTestService(tt.args.testRepository)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTestService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTestService() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTestServiceImpl_CreateTest(t *testing.T) {
	testutils.MonkeyAll()

	webService := &models.WebService{Id: 1}
	request := models.TestRequest{
		Path:        "/users",
		HttpMethod:  rshttp.MethodPost,
		ContentType: rshttp.MIMEApplicationJSON,
		RequestData: rsjson.MapJson{
			"name":   "sangil",
			"gender": "male",
		},
		Header:     nil,
		QueryParam: nil,
	}

	testWithoutId := &models.Test{
		DefaultValidateChecker: rsmodels.ValidatedDefaultValidateChecker,
		WebServiceId:           1,
		Path:                   request.Path,
		HttpMethod:             request.HttpMethod,
		ContentType:            request.ContentType,
		RequestData:            request.RequestData,
		Header:                 request.Header,
		QueryParam:             request.QueryParam,
		Created:                time.Now(),
		LastModified:           time.Now(),
	}

	test := &models.Test{
		DefaultValidateChecker: rsmodels.ValidatedDefaultValidateChecker,
		Id:                     1,
		WebServiceId:           1,
		Path:                   request.Path,
		HttpMethod:             request.HttpMethod,
		ContentType:            request.ContentType,
		RequestData:            request.RequestData,
		Header:                 request.Header,
		QueryParam:             request.QueryParam,
		Created:                time.Now(),
		LastModified:           time.Now(),
	}

	type args struct {
		webService *models.WebService
		request    models.TestRequest
	}
	tests := []struct {
		name     string
		args     args
		mockFunc testMockFunc
		want     *models.Test
		wantErr  *amerr.ErrorWithLanguage
	}{
		{
			name: "pass",
			args: args{
				webService: webService,
				request:    request,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockTestRepository *mocks.TestRepository) {
				mockTestRepository.
					On("Create", rsdb.GetConnection(), testWithoutId).
					Run(func(args mock.Arguments) {
						arg := args.Get(1).(*models.Test)
						arg.Id = 1
					}).
					Return(nil)
			},
			want:    test,
			wantErr: nil,
		},
		{
			name: "invalid argument",
			args: args{
				// webService: webService,
				request: request,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockTestRepository *mocks.TestRepository) {},
			want:     nil,
			wantErr:  amerr.GetErrInternalServer(),
		},
		{
			name: "invalid argument",
			args: args{
				webService: webService,
				// request: request,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockTestRepository *mocks.TestRepository) {},
			want:     nil,
			wantErr:  amerr.GetErrInternalServer(),
		},
		{
			name: "invalid request",
			args: args{
				webService: webService,
				request: models.TestRequest{
					Path: "/?a/a",
				},
			},
			mockFunc: func(mockConn *mocks2.Connection, mockTestRepository *mocks.TestRepository) {
				mockTestRepository.
					On("Create", rsdb.GetConnection(), testWithoutId).
					Run(func(args mock.Arguments) {
						arg := args.Get(1).(*models.Test)
						arg.Id = 1
					}).
					Return(nil)
			},
			want:    nil,
			wantErr: amerr.GetErrorsFromCode(amerr.ErrBadRequest),
		},
		{
			name: "WebService not found",
			args: args{
				webService: webService,
				request:    request,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockTestRepository *mocks.TestRepository) {
				mockTestRepository.
					On("Create", rsdb.GetConnection(), testWithoutId).
					Return(rsdb.ErrForeignKeyConstraint)
			},
			want:    nil,
			wantErr: amerr.GetErrorsFromCode(amerr.ErrWebServiceNotFound),
		},
		{
			name: "duplicated test",
			args: args{
				webService: webService,
				request:    request,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockTestRepository *mocks.TestRepository) {
				mockTestRepository.
					On("Create", rsdb.GetConnection(), testWithoutId).
					Return(rsdb.ErrDuplicateData)
			},
			want:    nil,
			wantErr: amerr.GetErrorsFromCode(amerr.ErrDuplicatedTest),
		},
		{
			name: "unexpected TestRepository.Create error",
			args: args{
				webService: webService,
				request:    request,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockTestRepository *mocks.TestRepository) {
				mockTestRepository.
					On("Create", rsdb.GetConnection(), testWithoutId).
					Return(rserrors.ErrUnexpected)
			},
			want:    nil,
			wantErr: amerr.GetErrInternalServer(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTestRepository := &mocks.TestRepository{}
			mockConn := &mocks2.Connection{}
			testutils.MonkeyGetConnection(mockConn)
			tt.mockFunc(mockConn, mockTestRepository)

			service := &TestServiceImpl{
				testRepository: mockTestRepository,
			}

			got, gotErr := service.CreateTest(tt.args.webService, tt.args.request)
			assert.Equal(t, tt.wantErr, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTestServiceImpl_GetTestById(t *testing.T) {
	testutils.MonkeyAll()

	testWithId := &models.Test{Id: 1}
	test := &models.Test{
		Id:           1,
		WebServiceId: 1,
		Path:         "/path/to/uri",
		HttpMethod:   rshttp.MethodGet,
		ContentType:  rshttp.MIMEApplicationJSON,
		RequestData:  nil,
		Header:       nil,
		QueryParam:   nil,
		Created:      time.Now(),
		LastModified: time.Now(),
	}

	type args struct {
		test *models.Test
	}
	tests := []struct {
		name     string
		args     args
		mockFunc testMockFunc
		want     *models.Test
		wantErr  *amerr.ErrorWithLanguage
	}{
		{
			name: "pass",
			args: args{
				test: testWithId,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockTestRepository *mocks.TestRepository) {
				mockTestRepository.On("GetById", rsdb.GetConnection(), testWithId).
					Run(func(args mock.Arguments) {
						arg := args.Get(1).(*models.Test)
						*arg = *test
					}).Return(nil)
			},
			want:    test,
			wantErr: nil,
		},
		{
			name: "test not found",
			args: args{
				test: testWithId,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockTestRepository *mocks.TestRepository) {
				mockTestRepository.On("GetById", rsdb.GetConnection(), testWithId).
					Return(rsdb.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: amerr.GetErrorsFromCode(amerr.ErrTestNotFound),
		},
		{
			name: "unexpected Repository.GetById error",
			args: args{
				test: testWithId,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockTestRepository *mocks.TestRepository) {
				mockTestRepository.On("GetById", rsdb.GetConnection(), testWithId).
					Return(rserrors.ErrUnexpected)
			},
			want:    nil,
			wantErr: amerr.GetErrInternalServer(),
		},
		{
			name: "invalid test",
			args: args{
				test: nil,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockTestRepository *mocks.TestRepository) {},
			want:     nil,
			wantErr:  amerr.GetErrInternalServer(),
		},
		{
			name: "invalid test",
			args: args{
				test: &models.Test{},
			},
			mockFunc: func(mockConn *mocks2.Connection, mockTestRepository *mocks.TestRepository) {},
			want:     nil,
			wantErr:  amerr.GetErrInternalServer(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTestRepository := &mocks.TestRepository{}
			mockConn := &mocks2.Connection{}
			testutils.MonkeyGetConnection(mockConn)
			tt.mockFunc(mockConn, mockTestRepository)

			service := &TestServiceImpl{
				testRepository: mockTestRepository,
			}
			assert.Equal(t, tt.wantErr, service.GetTestById(tt.args.test))
			if tt.wantErr == nil {
				assert.Equal(t, tt.want, tt.args.test)
			}
		})
	}
}

func TestTestServiceImpl_DeleteTestById(t *testing.T) {
	testutils.MonkeyAll()

	testWithId := &models.Test{Id: 1}

	type args struct {
		test *models.Test
	}
	tests := []struct {
		name     string
		args     args
		mockFunc testMockFunc
		wantErr  *amerr.ErrorWithLanguage
	}{
		{
			name: "pass",
			args: args{
				test: testWithId,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockTestRepository *mocks.TestRepository) {
				mockTestRepository.On("DeleteById", rsdb.GetConnection(), testWithId).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "invalid parameter",
			args: args{
				test: nil,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockTestRepository *mocks.TestRepository) {
			},
			wantErr: amerr.GetErrInternalServer(),
		},
		{
			name: "invalid parameter",
			args: args{
				test: &models.Test{},
			},
			mockFunc: func(mockConn *mocks2.Connection, mockTestRepository *mocks.TestRepository) {
			},
			wantErr: amerr.GetErrInternalServer(),
		},
		{
			name: "unexpected TestRepository error",
			args: args{
				test: testWithId,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockTestRepository *mocks.TestRepository) {
				mockTestRepository.On("DeleteById", rsdb.GetConnection(), testWithId).
					Return(rserrors.ErrUnexpected)
			},
			wantErr: amerr.GetErrInternalServer(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTestRepository := &mocks.TestRepository{}
			mockConn := &mocks2.Connection{}
			testutils.MonkeyGetConnection(mockConn)
			tt.mockFunc(mockConn, mockTestRepository)

			service := &TestServiceImpl{
				testRepository: mockTestRepository,
			}

			assert.Equal(t, tt.wantErr, service.DeleteTestById(tt.args.test))
		})
	}
}
