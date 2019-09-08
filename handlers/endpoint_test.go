package handlers

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/realsangil/apimonitor/middlewares/mocks"
	"github.com/realsangil/apimonitor/models"
	"github.com/realsangil/apimonitor/pkg/amerr"
	"github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/rshttp"
	"github.com/realsangil/apimonitor/pkg/testutils"
	"github.com/realsangil/apimonitor/services"
	mocks2 "github.com/realsangil/apimonitor/services/mocks"
)

const langForTest = "ko"

func setContextLanguageForTest(ctx *mocks.Context) {
	ctx.On("Language").Return(langForTest)
}

type endpointMockFunc func(mockContext *mocks.Context, mockEndpointService *mocks2.EndpointService, mockWebServiceService *mocks2.WebServiceService)

func TestEndpointHandlerImpl_CreateEndpoint(t *testing.T) {
	testutils.MonkeyAll()

	webServiceId := int64(1)
	zeroInt64 := int64(0)

	webServiceWithId := &models.WebService{Id: webServiceId}
	webService := &models.WebService{
		DefaultValidateChecker: models.ValidatedDefaultValidateChecker,
		Id:                     1,
		Host:                   "realsangil.github.io",
		HttpSchema:             "https",
		Desc:                   "sangil's dev blog",
		Favicon:                "",
		Created:                time.Now(),
		LastModified:           time.Now(),
	}

	request := models.EndpointRequest{
		Path:        "/path/to/uri",
		HttpMethod:  rshttp.MethodGet,
		ContentType: rshttp.MIMEApplicationJSON,
		RequestData: nil,
		Header:      nil,
		QueryParam:  nil,
	}

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

	tests := []struct {
		name     string
		mockFunc endpointMockFunc
		wantErr  error
	}{
		{
			name: "pass",
			mockFunc: func(mockContext *mocks.Context, mockEndpointService *mocks2.EndpointService, mockWebServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", WebServiceIdParam).Return(webServiceId, nil)

				mockWebServiceService.On("GetWebServiceById", webServiceWithId).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*models.WebService)
						*arg = *webService
					}).
					Return(nil)

				mockContext.On("Bind", &models.EndpointRequest{}).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*models.EndpointRequest)
						*arg = request
					}).
					Return(nil)

				mockEndpointService.On("CreateEndpoint", webService, request).
					Return(endpoint, nil)

				mockContext.On("JSON", http.StatusOK, endpoint).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "invalid webServiceId",
			mockFunc: func(mockContext *mocks.Context, mockEndpointService *mocks2.EndpointService, mockWebServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", WebServiceIdParam).Return(zeroInt64, rserrors.ErrInvalidParameter)
			},
			wantErr: amerr.GetErrorsFromCode(amerr.ErrWebServiceNotFound).GetErrFromLanguage(langForTest),
		},
		{
			name: "unexpected webServiceService.GetWebServiceById error",
			mockFunc: func(mockContext *mocks.Context, mockEndpointService *mocks2.EndpointService, mockWebServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", WebServiceIdParam).Return(webServiceId, nil)

				mockWebServiceService.On("GetWebServiceById", webServiceWithId).
					Return(amerr.GetErrInternalServer())
			},
			wantErr: amerr.GetErrInternalServer().GetErrFromLanguage(langForTest),
		},
		{
			name: "bind error",
			mockFunc: func(mockContext *mocks.Context, mockEndpointService *mocks2.EndpointService, mockWebServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", WebServiceIdParam).Return(webServiceId, nil)

				mockWebServiceService.On("GetWebServiceById", webServiceWithId).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*models.WebService)
						*arg = *webService
					}).
					Return(nil)

				mockContext.On("Bind", &models.EndpointRequest{}).
					Return(rserrors.ErrUnexpected)
			},
			wantErr: amerr.GetErrorsFromCode(amerr.ErrBadRequest).GetErrFromLanguage(langForTest),
		},
		{
			name: "unexpected EndpointService.CreateEndpoint error",
			mockFunc: func(mockContext *mocks.Context, mockEndpointService *mocks2.EndpointService, mockWebServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", WebServiceIdParam).Return(webServiceId, nil)

				mockWebServiceService.On("GetWebServiceById", webServiceWithId).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*models.WebService)
						*arg = *webService
					}).
					Return(nil)

				mockContext.On("Bind", &models.EndpointRequest{}).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*models.EndpointRequest)
						*arg = request
					}).
					Return(nil)

				mockEndpointService.On("CreateEndpoint", webService, request).
					Return(nil, amerr.GetErrInternalServer())
			},
			wantErr: amerr.GetErrInternalServer().GetErrFromLanguage(langForTest),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContext := &mocks.Context{}
			mockEndpointService := &mocks2.EndpointService{}
			mockWebServiceService := &mocks2.WebServiceService{}

			tt.mockFunc(mockContext, mockEndpointService, mockWebServiceService)
			handler := &EndpointHandlerImpl{
				webServiceService: mockWebServiceService,
				endpointService:   mockEndpointService,
			}
			assert.Equal(t, tt.wantErr, handler.CreateEndpoint(mockContext))
		})
	}
}

func TestEndpointHandlerImpl_GetEndpoint(t *testing.T) {
	testutils.MonkeyAll()

	zeroInt64 := int64(0)
	endpointId := int64(1)

	endpointWithId := &models.Endpoint{Id: endpointId}
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

	tests := []struct {
		name     string
		mockFunc endpointMockFunc
		wantErr  error
	}{
		{
			name: "pass",
			mockFunc: func(mockContext *mocks.Context, mockEndpointService *mocks2.EndpointService, mockWebServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", EndpointIdParam).Return(endpointId, nil)

				mockEndpointService.On("GetEndpointById", endpointWithId).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*models.Endpoint)
						*arg = *endpoint
					}).
					Return(nil)

				mockContext.On("JSON", http.StatusOK, endpoint).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "invalid endpoint id parameter",
			mockFunc: func(mockContext *mocks.Context, mockEndpointService *mocks2.EndpointService, webServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", EndpointIdParam).Return(zeroInt64, rserrors.ErrInvalidParameter)
			},
			wantErr: amerr.GetErrorsFromCode(amerr.ErrEndpointNotFound).GetErrFromLanguage(langForTest),
		},
		{
			name: "unexpected EndpointService.GetEndpointById error",
			mockFunc: func(mockContext *mocks.Context, mockEndpointService *mocks2.EndpointService, webServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", EndpointIdParam).Return(endpointId, nil)

				mockEndpointService.On("GetEndpointById", endpointWithId).
					Return(amerr.GetErrInternalServer())
			},
			wantErr: amerr.GetErrInternalServer().GetErrFromLanguage(langForTest),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContext := &mocks.Context{}
			mockEndpointService := &mocks2.EndpointService{}
			mockWebServiceService := &mocks2.WebServiceService{}

			tt.mockFunc(mockContext, mockEndpointService, mockWebServiceService)
			handler := &EndpointHandlerImpl{
				webServiceService: mockWebServiceService,
				endpointService:   mockEndpointService,
			}
			assert.Equal(t, tt.wantErr, handler.GetEndpoint(mockContext))
		})
	}
}

func TestEndpointHandlerImpl_DeleteEndpoint(t *testing.T) {
	testutils.MonkeyAll()

	zeroInt64 := int64(0)
	endpointId := int64(1)

	endpointWithId := &models.Endpoint{Id: endpointId}

	tests := []struct {
		name     string
		mockFunc endpointMockFunc
		wantErr  error
	}{
		{
			name: "pass",
			mockFunc: func(mockContext *mocks.Context, mockEndpointService *mocks2.EndpointService, mockWebServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", EndpointIdParam).Return(endpointId, nil)

				mockEndpointService.On("DeleteEndpointById", endpointWithId).
					Return(nil)

				mockContext.On("JSON", http.StatusOK, nil).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "invalid endpoint id parameter",
			mockFunc: func(mockContext *mocks.Context, mockEndpointService *mocks2.EndpointService, webServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", EndpointIdParam).Return(zeroInt64, rserrors.ErrInvalidParameter)
			},
			wantErr: amerr.GetErrorsFromCode(amerr.ErrEndpointNotFound).GetErrFromLanguage(langForTest),
		},
		{
			name: "unexpected EndpointService.GetEndpointById error",
			mockFunc: func(mockContext *mocks.Context, mockEndpointService *mocks2.EndpointService, webServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", EndpointIdParam).Return(endpointId, nil)

				mockEndpointService.On("DeleteEndpointById", endpointWithId).
					Return(amerr.GetErrInternalServer())
			},
			wantErr: amerr.GetErrInternalServer().GetErrFromLanguage(langForTest),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContext := &mocks.Context{}
			mockEndpointService := &mocks2.EndpointService{}
			mockWebServiceService := &mocks2.WebServiceService{}

			tt.mockFunc(mockContext, mockEndpointService, mockWebServiceService)
			handler := &EndpointHandlerImpl{
				webServiceService: mockWebServiceService,
				endpointService:   mockEndpointService,
			}
			assert.Equal(t, tt.wantErr, handler.DeleteEndpoint(mockContext))
		})
	}
}

func TestNewEndpointHandler(t *testing.T) {
	type args struct {
		webServiceService services.WebServiceService
		endpointService   services.EndpointService
	}
	tests := []struct {
		name    string
		args    args
		want    EndpointHandler
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				webServiceService: &mocks2.WebServiceService{},
				endpointService:   &mocks2.EndpointService{},
			},
			want: &EndpointHandlerImpl{
				webServiceService: &mocks2.WebServiceService{},
				endpointService:   &mocks2.EndpointService{},
			},
			wantErr: false,
		},
		{
			name: "invalid field",
			args: args{
				// webServiceService: &mocks2.WebServiceService{},
				endpointService: &mocks2.EndpointService{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid field",
			args: args{
				webServiceService: &mocks2.WebServiceService{},
				// endpointService:   &mocks2.EndpointService{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEndpointHandler(tt.args.webServiceService, tt.args.endpointService)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEndpointHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEndpointHandler() got = %v, want %v", got, tt.want)
			}
		})
	}
}
