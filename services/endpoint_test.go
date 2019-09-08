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

type endpointMockFunc func(mockConn *mocks2.Connection, mockEndpointRepository *mocks.EndpointRepository)

func TestNewEndpointService(t *testing.T) {
	type args struct {
		endpointRepository repositories.EndpointRepository
	}
	tests := []struct {
		name    string
		args    args
		want    EndpointService
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				endpointRepository: &mocks.EndpointRepository{},
			},
			want: &EndpointServiceImpl{
				endpointRepository: &mocks.EndpointRepository{},
			},
			wantErr: false,
		},
		{
			name: "invalid parameter",
			args: args{
				// endpointRepository: &mocks.EndpointRepository{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEndpointService(tt.args.endpointRepository)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEndpointService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEndpointService() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndpointServiceImpl_CreateEndpoint(t *testing.T) {
	testutils.MonkeyAll()

	webService := &models.WebService{Id: 1}
	request := models.EndpointRequest{
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

	endpointWitoutId := &models.Endpoint{
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

	endpoint := &models.Endpoint{
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
		request    models.EndpointRequest
	}
	tests := []struct {
		name     string
		args     args
		mockFunc endpointMockFunc
		want     *models.Endpoint
		wantErr  *amerr.ErrorWithLanguage
	}{
		{
			name: "pass",
			args: args{
				webService: webService,
				request:    request,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockEndpointRepository *mocks.EndpointRepository) {
				mockEndpointRepository.
					On("Create", rsdb.GetConnection(), endpointWitoutId).
					Run(func(args mock.Arguments) {
						arg := args.Get(1).(*models.Endpoint)
						arg.Id = 1
					}).
					Return(nil)
			},
			want:    endpoint,
			wantErr: nil,
		},
		{
			name: "invalid argument",
			args: args{
				// webService: webService,
				request: request,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockEndpointRepository *mocks.EndpointRepository) {},
			want:     nil,
			wantErr:  amerr.GetErrInternalServer(),
		},
		{
			name: "invalid argument",
			args: args{
				webService: webService,
				// request: request,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockEndpointRepository *mocks.EndpointRepository) {},
			want:     nil,
			wantErr:  amerr.GetErrInternalServer(),
		},
		{
			name: "invalid request",
			args: args{
				webService: webService,
				request: models.EndpointRequest{
					Path: "/?a/a",
				},
			},
			mockFunc: func(mockConn *mocks2.Connection, mockEndpointRepository *mocks.EndpointRepository) {
				mockEndpointRepository.
					On("Create", rsdb.GetConnection(), endpointWitoutId).
					Run(func(args mock.Arguments) {
						arg := args.Get(1).(*models.Endpoint)
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
			mockFunc: func(mockConn *mocks2.Connection, mockEndpointRepository *mocks.EndpointRepository) {
				mockEndpointRepository.
					On("Create", rsdb.GetConnection(), endpointWitoutId).
					Return(rsdb.ErrForeignKeyConstraint)
			},
			want:    nil,
			wantErr: amerr.GetErrorsFromCode(amerr.ErrWebServiceNotFound),
		},
		{
			name: "duplicated endpoint",
			args: args{
				webService: webService,
				request:    request,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockEndpointRepository *mocks.EndpointRepository) {
				mockEndpointRepository.
					On("Create", rsdb.GetConnection(), endpointWitoutId).
					Return(rsdb.ErrDuplicateData)
			},
			want:    nil,
			wantErr: amerr.GetErrorsFromCode(amerr.ErrDuplicatedEndpoint),
		},
		{
			name: "unexpected EndpointRepository.Create error",
			args: args{
				webService: webService,
				request:    request,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockEndpointRepository *mocks.EndpointRepository) {
				mockEndpointRepository.
					On("Create", rsdb.GetConnection(), endpointWitoutId).
					Return(rserrors.ErrUnexpected)
			},
			want:    nil,
			wantErr: amerr.GetErrInternalServer(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockEndpointRepository := &mocks.EndpointRepository{}
			mockConn := &mocks2.Connection{}
			testutils.MonkeyGetConnection(mockConn)
			tt.mockFunc(mockConn, mockEndpointRepository)

			service := &EndpointServiceImpl{
				endpointRepository: mockEndpointRepository,
			}

			got, gotErr := service.CreateEndpoint(tt.args.webService, tt.args.request)
			assert.Equal(t, tt.wantErr, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEndpointServiceImpl_GetEndpointById(t *testing.T) {
	testutils.MonkeyAll()

	endpointWithId := &models.Endpoint{Id: 1}
	endpoint := &models.Endpoint{
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
		endpoint *models.Endpoint
	}
	tests := []struct {
		name     string
		args     args
		mockFunc endpointMockFunc
		want     *models.Endpoint
		wantErr  *amerr.ErrorWithLanguage
	}{
		{
			name: "pass",
			args: args{
				endpoint: endpointWithId,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockEndpointRepository *mocks.EndpointRepository) {
				mockEndpointRepository.On("GetById", rsdb.GetConnection(), endpointWithId).
					Run(func(args mock.Arguments) {
						arg := args.Get(1).(*models.Endpoint)
						*arg = *endpoint
					}).Return(nil)
			},
			want:    endpoint,
			wantErr: nil,
		},
		{
			name: "endpoint not found",
			args: args{
				endpoint: endpointWithId,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockEndpointRepository *mocks.EndpointRepository) {
				mockEndpointRepository.On("GetById", rsdb.GetConnection(), endpointWithId).
					Return(rsdb.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: amerr.GetErrorsFromCode(amerr.ErrEndpointNotFound),
		},
		{
			name: "unexpected Repository.GetById error",
			args: args{
				endpoint: endpointWithId,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockEndpointRepository *mocks.EndpointRepository) {
				mockEndpointRepository.On("GetById", rsdb.GetConnection(), endpointWithId).
					Return(rserrors.ErrUnexpected)
			},
			want:    nil,
			wantErr: amerr.GetErrInternalServer(),
		},
		{
			name: "invalid endpoint",
			args: args{
				endpoint: nil,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockEndpointRepository *mocks.EndpointRepository) {},
			want:     nil,
			wantErr:  amerr.GetErrInternalServer(),
		},
		{
			name: "invalid endpoint",
			args: args{
				endpoint: &models.Endpoint{},
			},
			mockFunc: func(mockConn *mocks2.Connection, mockEndpointRepository *mocks.EndpointRepository) {},
			want:     nil,
			wantErr:  amerr.GetErrInternalServer(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockEndpointRepository := &mocks.EndpointRepository{}
			mockConn := &mocks2.Connection{}
			testutils.MonkeyGetConnection(mockConn)
			tt.mockFunc(mockConn, mockEndpointRepository)

			service := &EndpointServiceImpl{
				endpointRepository: mockEndpointRepository,
			}
			assert.Equal(t, tt.wantErr, service.GetEndpointById(tt.args.endpoint))
			if tt.wantErr == nil {
				assert.Equal(t, tt.want, tt.args.endpoint)
			}
		})
	}
}

func TestEndpointServiceImpl_DeleteEndpointById(t *testing.T) {
	testutils.MonkeyAll()

	endpointWithId := &models.Endpoint{Id: 1}

	type args struct {
		endpoint *models.Endpoint
	}
	tests := []struct {
		name     string
		args     args
		mockFunc endpointMockFunc
		wantErr  *amerr.ErrorWithLanguage
	}{
		{
			name: "pass",
			args: args{
				endpoint: endpointWithId,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockEndpointRepository *mocks.EndpointRepository) {
				mockEndpointRepository.On("DeleteById", rsdb.GetConnection(), endpointWithId).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "invalid parameter",
			args: args{
				endpoint: nil,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockEndpointRepository *mocks.EndpointRepository) {
			},
			wantErr: amerr.GetErrInternalServer(),
		},
		{
			name: "invalid parameter",
			args: args{
				endpoint: &models.Endpoint{},
			},
			mockFunc: func(mockConn *mocks2.Connection, mockEndpointRepository *mocks.EndpointRepository) {
			},
			wantErr: amerr.GetErrInternalServer(),
		},
		{
			name: "unexpected EndpointRepository error",
			args: args{
				endpoint: endpointWithId,
			},
			mockFunc: func(mockConn *mocks2.Connection, mockEndpointRepository *mocks.EndpointRepository) {
				mockEndpointRepository.On("DeleteById", rsdb.GetConnection(), endpointWithId).
					Return(rserrors.ErrUnexpected)
			},
			wantErr: amerr.GetErrInternalServer(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockEndpointRepository := &mocks.EndpointRepository{}
			mockConn := &mocks2.Connection{}
			testutils.MonkeyGetConnection(mockConn)
			tt.mockFunc(mockConn, mockEndpointRepository)

			service := &EndpointServiceImpl{
				endpointRepository: mockEndpointRepository,
			}

			assert.Equal(t, tt.wantErr, service.DeleteEndpointById(tt.args.endpoint))
		})
	}
}
