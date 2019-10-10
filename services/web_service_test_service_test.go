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
	"github.com/realsangil/apimonitor/pkg/testutils"
	"github.com/realsangil/apimonitor/repositories"
	"github.com/realsangil/apimonitor/repositories/mocks"
)

type webServiceTestMockFunc func(mockConn *mocks2.Connection, mockWebServiceTestRepository *mocks.WebServiceTestRepository)

func TestNewWebServiceTestService(t *testing.T) {
	type args struct {
		webServiceTestRepository repositories.WebServiceTestRepository
	}
	tests := []struct {
		name    string
		args    args
		want    WebServiceTestService
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				webServiceTestRepository: &mocks.WebServiceTestRepository{},
			},
			want: &WebServiceTestServiceImpl{
				webServiceTestRepository: &mocks.WebServiceTestRepository{},
			},
			wantErr: false,
		},
		{
			name: "invalid parameter",
			args: args{
				// webServiceTestRepository: &mocks.WebServiceTestRepository{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewWebServiceTestService(tt.args.webServiceTestRepository)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWebServiceTestService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWebServiceTestService() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebServiceTestServiceImpl_CreateWebServiceTest(t *testing.T) {
	testutils.MonkeyAll()

	webService := &models.WebService{Id: 1}
	request := models.WebServiceTestRequest{
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

	webServiceTestWithoutId := &models.WebServiceTest{
		DefaultValidateChecker: models.ValidatedDefaultValidateChecker,
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

	webServiceTest := &models.WebServiceTest{
		DefaultValidateChecker: models.ValidatedDefaultValidateChecker,
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
		request    models.WebServiceTestRequest
	}
	tests := []struct {
		name     string
		args     args
		mockFunc webServiceTestMockFunc
		want     *models.WebServiceTest
		wantErr  *amerr.ErrorWithLanguage
	}{
		{
			name: "pass",
			args: args{
				webService: webService,
				request:    request,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockWebServiceTestRepository *mocks.WebServiceTestRepository) {
				mockWebServiceTestRepository.
					On("Create", rsdb.GetConnection(), webServiceTestWithoutId).
					Run(func(args mock.Arguments) {
						arg := args.Get(1).(*models.WebServiceTest)
						arg.Id = 1
					}).
					Return(nil)
			},
			want:    webServiceTest,
			wantErr: nil,
		},
		{
			name: "invalid argument",
			args: args{
				// webService: webService,
				request: request,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockWebServiceTestRepository *mocks.WebServiceTestRepository) {},
			want:     nil,
			wantErr:  amerr.GetErrInternalServer(),
		},
		{
			name: "invalid argument",
			args: args{
				webService: webService,
				// request: request,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockWebServiceTestRepository *mocks.WebServiceTestRepository) {},
			want:     nil,
			wantErr:  amerr.GetErrInternalServer(),
		},
		{
			name: "invalid request",
			args: args{
				webService: webService,
				request: models.WebServiceTestRequest{
					Path: "/?a/a",
				},
			},
			mockFunc: func(mockConn *mocks2.Connection, mockWebServiceTestRepository *mocks.WebServiceTestRepository) {
				mockWebServiceTestRepository.
					On("Create", rsdb.GetConnection(), webServiceTestWithoutId).
					Run(func(args mock.Arguments) {
						arg := args.Get(1).(*models.WebServiceTest)
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
			mockFunc: func(mockConn *mocks2.Connection, mockWebServiceTestRepository *mocks.WebServiceTestRepository) {
				mockWebServiceTestRepository.
					On("Create", rsdb.GetConnection(), webServiceTestWithoutId).
					Return(rsdb.ErrForeignKeyConstraint)
			},
			want:    nil,
			wantErr: amerr.GetErrorsFromCode(amerr.ErrWebServiceNotFound),
		},
		{
			name: "duplicated webServiceTest",
			args: args{
				webService: webService,
				request:    request,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockWebServiceTestRepository *mocks.WebServiceTestRepository) {
				mockWebServiceTestRepository.
					On("Create", rsdb.GetConnection(), webServiceTestWithoutId).
					Return(rsdb.ErrDuplicateData)
			},
			want:    nil,
			wantErr: amerr.GetErrorsFromCode(amerr.ErrDuplicatedWebServiceTest),
		},
		{
			name: "unexpected WebServiceTestRepository.Create error",
			args: args{
				webService: webService,
				request:    request,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockWebServiceTestRepository *mocks.WebServiceTestRepository) {
				mockWebServiceTestRepository.
					On("Create", rsdb.GetConnection(), webServiceTestWithoutId).
					Return(rserrors.ErrUnexpected)
			},
			want:    nil,
			wantErr: amerr.GetErrInternalServer(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWebServiceTestRepository := &mocks.WebServiceTestRepository{}
			mockConn := &mocks2.Connection{}
			testutils.MonkeyGetConnection(mockConn)
			tt.mockFunc(mockConn, mockWebServiceTestRepository)

			service := &WebServiceTestServiceImpl{
				webServiceTestRepository: mockWebServiceTestRepository,
			}

			got, gotErr := service.CreateWebServiceTest(tt.args.webService, tt.args.request)
			assert.Equal(t, tt.wantErr, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWebServiceTestServiceImpl_GetWebServiceTestById(t *testing.T) {
	testutils.MonkeyAll()

	webServiceTestWithId := &models.WebServiceTest{Id: 1}
	webServiceTest := &models.WebServiceTest{
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
		webServiceTest *models.WebServiceTest
	}
	tests := []struct {
		name     string
		args     args
		mockFunc webServiceTestMockFunc
		want     *models.WebServiceTest
		wantErr  *amerr.ErrorWithLanguage
	}{
		{
			name: "pass",
			args: args{
				webServiceTest: webServiceTestWithId,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockWebServiceTestRepository *mocks.WebServiceTestRepository) {
				mockWebServiceTestRepository.On("GetById", rsdb.GetConnection(), webServiceTestWithId).
					Run(func(args mock.Arguments) {
						arg := args.Get(1).(*models.WebServiceTest)
						*arg = *webServiceTest
					}).Return(nil)
			},
			want:    webServiceTest,
			wantErr: nil,
		},
		{
			name: "webServiceTest not found",
			args: args{
				webServiceTest: webServiceTestWithId,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockWebServiceTestRepository *mocks.WebServiceTestRepository) {
				mockWebServiceTestRepository.On("GetById", rsdb.GetConnection(), webServiceTestWithId).
					Return(rsdb.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: amerr.GetErrorsFromCode(amerr.ErrWebServiceTestNotFound),
		},
		{
			name: "unexpected Repository.GetById error",
			args: args{
				webServiceTest: webServiceTestWithId,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockWebServiceTestRepository *mocks.WebServiceTestRepository) {
				mockWebServiceTestRepository.On("GetById", rsdb.GetConnection(), webServiceTestWithId).
					Return(rserrors.ErrUnexpected)
			},
			want:    nil,
			wantErr: amerr.GetErrInternalServer(),
		},
		{
			name: "invalid webServiceTest",
			args: args{
				webServiceTest: nil,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockWebServiceTestRepository *mocks.WebServiceTestRepository) {},
			want:     nil,
			wantErr:  amerr.GetErrInternalServer(),
		},
		{
			name: "invalid webServiceTest",
			args: args{
				webServiceTest: &models.WebServiceTest{},
			},
			mockFunc: func(mockConn *mocks2.Connection, mockWebServiceTestRepository *mocks.WebServiceTestRepository) {},
			want:     nil,
			wantErr:  amerr.GetErrInternalServer(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWebServiceTestRepository := &mocks.WebServiceTestRepository{}
			mockConn := &mocks2.Connection{}
			testutils.MonkeyGetConnection(mockConn)
			tt.mockFunc(mockConn, mockWebServiceTestRepository)

			service := &WebServiceTestServiceImpl{
				webServiceTestRepository: mockWebServiceTestRepository,
			}
			assert.Equal(t, tt.wantErr, service.GetWebServiceTestById(tt.args.webServiceTest))
			if tt.wantErr == nil {
				assert.Equal(t, tt.want, tt.args.webServiceTest)
			}
		})
	}
}

func TestWebServiceTestServiceImpl_DeleteWebServiceTestById(t *testing.T) {
	testutils.MonkeyAll()

	webServiceTestWithId := &models.WebServiceTest{Id: 1}

	type args struct {
		webServiceTest *models.WebServiceTest
	}
	tests := []struct {
		name     string
		args     args
		mockFunc webServiceTestMockFunc
		wantErr  *amerr.ErrorWithLanguage
	}{
		{
			name: "pass",
			args: args{
				webServiceTest: webServiceTestWithId,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockWebServiceTestRepository *mocks.WebServiceTestRepository) {
				mockWebServiceTestRepository.On("DeleteById", rsdb.GetConnection(), webServiceTestWithId).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "invalid parameter",
			args: args{
				webServiceTest: nil,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockWebServiceTestRepository *mocks.WebServiceTestRepository) {
			},
			wantErr: amerr.GetErrInternalServer(),
		},
		{
			name: "invalid parameter",
			args: args{
				webServiceTest: &models.WebServiceTest{},
			},
			mockFunc: func(mockConn *mocks2.Connection, mockWebServiceTestRepository *mocks.WebServiceTestRepository) {
			},
			wantErr: amerr.GetErrInternalServer(),
		},
		{
			name: "unexpected WebServiceTestRepository error",
			args: args{
				webServiceTest: webServiceTestWithId,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockWebServiceTestRepository *mocks.WebServiceTestRepository) {
				mockWebServiceTestRepository.On("DeleteById", rsdb.GetConnection(), webServiceTestWithId).
					Return(rserrors.ErrUnexpected)
			},
			wantErr: amerr.GetErrInternalServer(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWebServiceTestRepository := &mocks.WebServiceTestRepository{}
			mockConn := &mocks2.Connection{}
			testutils.MonkeyGetConnection(mockConn)
			tt.mockFunc(mockConn, mockWebServiceTestRepository)

			service := &WebServiceTestServiceImpl{
				webServiceTestRepository: mockWebServiceTestRepository,
			}

			assert.Equal(t, tt.wantErr, service.DeleteWebServiceTestById(tt.args.webServiceTest))
		})
	}
}
