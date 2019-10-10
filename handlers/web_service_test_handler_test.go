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

type webServiceTestMockFunc func(mockContext *mocks.Context, mockWebServiceTestService *mocks2.WebServiceTestService, mockWebServiceService *mocks2.WebServiceService)

func TestWebServiceTestHandlerImpl_CreateWebServiceTest(t *testing.T) {
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

	request := models.WebServiceTestRequest{
		Path:        "/path/to/uri",
		HttpMethod:  rshttp.MethodGet,
		ContentType: rshttp.MIMEApplicationJSON,
		RequestData: nil,
		Header:      nil,
		QueryParam:  nil,
	}

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

	tests := []struct {
		name     string
		mockFunc webServiceTestMockFunc
		wantErr  error
	}{
		{
			name: "pass",
			mockFunc: func(mockContext *mocks.Context, mockWebServiceTestService *mocks2.WebServiceTestService, mockWebServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", WebServiceIdParam).Return(webServiceId, nil)

				mockWebServiceService.On("GetWebServiceById", webServiceWithId).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*models.WebService)
						*arg = *webService
					}).
					Return(nil)

				mockContext.On("Bind", &models.WebServiceTestRequest{}).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*models.WebServiceTestRequest)
						*arg = request
					}).
					Return(nil)

				mockWebServiceTestService.On("CreateWebServiceTest", webService, request).
					Return(webServiceTest, nil)

				mockContext.On("JSON", http.StatusOK, webServiceTest).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "invalid webServiceId",
			mockFunc: func(mockContext *mocks.Context, mockWebServiceTestService *mocks2.WebServiceTestService, mockWebServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", WebServiceIdParam).Return(zeroInt64, rserrors.ErrInvalidParameter)
			},
			wantErr: amerr.GetErrorsFromCode(amerr.ErrWebServiceNotFound).GetErrFromLanguage(langForTest),
		},
		{
			name: "unexpected webServiceService.GetWebServiceById error",
			mockFunc: func(mockContext *mocks.Context, mockWebServiceTestService *mocks2.WebServiceTestService, mockWebServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", WebServiceIdParam).Return(webServiceId, nil)

				mockWebServiceService.On("GetWebServiceById", webServiceWithId).
					Return(amerr.GetErrInternalServer())
			},
			wantErr: amerr.GetErrInternalServer().GetErrFromLanguage(langForTest),
		},
		{
			name: "bind error",
			mockFunc: func(mockContext *mocks.Context, mockWebServiceTestService *mocks2.WebServiceTestService, mockWebServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", WebServiceIdParam).Return(webServiceId, nil)

				mockWebServiceService.On("GetWebServiceById", webServiceWithId).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*models.WebService)
						*arg = *webService
					}).
					Return(nil)

				mockContext.On("Bind", &models.WebServiceTestRequest{}).
					Return(rserrors.ErrUnexpected)
			},
			wantErr: amerr.GetErrorsFromCode(amerr.ErrBadRequest).GetErrFromLanguage(langForTest),
		},
		{
			name: "unexpected WebServiceTestService.CreateWebServiceTest error",
			mockFunc: func(mockContext *mocks.Context, mockWebServiceTestService *mocks2.WebServiceTestService, mockWebServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", WebServiceIdParam).Return(webServiceId, nil)

				mockWebServiceService.On("GetWebServiceById", webServiceWithId).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*models.WebService)
						*arg = *webService
					}).
					Return(nil)

				mockContext.On("Bind", &models.WebServiceTestRequest{}).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*models.WebServiceTestRequest)
						*arg = request
					}).
					Return(nil)

				mockWebServiceTestService.On("CreateWebServiceTest", webService, request).
					Return(nil, amerr.GetErrInternalServer())
			},
			wantErr: amerr.GetErrInternalServer().GetErrFromLanguage(langForTest),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContext := &mocks.Context{}
			mockWebServiceTestService := &mocks2.WebServiceTestService{}
			mockWebServiceService := &mocks2.WebServiceService{}

			tt.mockFunc(mockContext, mockWebServiceTestService, mockWebServiceService)
			handler := &WebServiceTestHandlerImpl{
				webServiceService:     mockWebServiceService,
				webServiceTestService: mockWebServiceTestService,
			}
			assert.Equal(t, tt.wantErr, handler.CreateWebServiceTest(mockContext))
		})
	}
}

func TestWebServiceTestHandlerImpl_GetWebServiceTest(t *testing.T) {
	testutils.MonkeyAll()

	zeroInt64 := int64(0)
	webServiceTestId := int64(1)

	webServiceTestWithId := &models.WebServiceTest{Id: webServiceTestId}
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

	tests := []struct {
		name     string
		mockFunc webServiceTestMockFunc
		wantErr  error
	}{
		{
			name: "pass",
			mockFunc: func(mockContext *mocks.Context, mockWebServiceTestService *mocks2.WebServiceTestService, mockWebServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", WebServiceTestIdParam).Return(webServiceTestId, nil)

				mockWebServiceTestService.On("GetWebServiceTestById", webServiceTestWithId).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*models.WebServiceTest)
						*arg = *webServiceTest
					}).
					Return(nil)

				mockContext.On("JSON", http.StatusOK, webServiceTest).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "invalid webServiceTest id parameter",
			mockFunc: func(mockContext *mocks.Context, mockWebServiceTestService *mocks2.WebServiceTestService, webServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", WebServiceTestIdParam).Return(zeroInt64, rserrors.ErrInvalidParameter)
			},
			wantErr: amerr.GetErrorsFromCode(amerr.ErrWebServiceTestNotFound).GetErrFromLanguage(langForTest),
		},
		{
			name: "unexpected WebServiceTestService.GetWebServiceTestById error",
			mockFunc: func(mockContext *mocks.Context, mockWebServiceTestService *mocks2.WebServiceTestService, webServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", WebServiceTestIdParam).Return(webServiceTestId, nil)

				mockWebServiceTestService.On("GetWebServiceTestById", webServiceTestWithId).
					Return(amerr.GetErrInternalServer())
			},
			wantErr: amerr.GetErrInternalServer().GetErrFromLanguage(langForTest),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContext := &mocks.Context{}
			mockWebServiceTestService := &mocks2.WebServiceTestService{}
			mockWebServiceService := &mocks2.WebServiceService{}

			tt.mockFunc(mockContext, mockWebServiceTestService, mockWebServiceService)
			handler := &WebServiceTestHandlerImpl{
				webServiceService:     mockWebServiceService,
				webServiceTestService: mockWebServiceTestService,
			}
			assert.Equal(t, tt.wantErr, handler.GetWebServiceTest(mockContext))
		})
	}
}

func TestWebServiceTestHandlerImpl_DeleteWebServiceTest(t *testing.T) {
	testutils.MonkeyAll()

	zeroInt64 := int64(0)
	webServiceTestId := int64(1)

	webServiceTestWithId := &models.WebServiceTest{Id: webServiceTestId}

	tests := []struct {
		name     string
		mockFunc webServiceTestMockFunc
		wantErr  error
	}{
		{
			name: "pass",
			mockFunc: func(mockContext *mocks.Context, mockWebServiceTestService *mocks2.WebServiceTestService, mockWebServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", WebServiceTestIdParam).Return(webServiceTestId, nil)

				mockWebServiceTestService.On("DeleteWebServiceTestById", webServiceTestWithId).
					Return(nil)

				mockContext.On("JSON", http.StatusOK, nil).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "invalid webServiceTest id parameter",
			mockFunc: func(mockContext *mocks.Context, mockWebServiceTestService *mocks2.WebServiceTestService, webServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", WebServiceTestIdParam).Return(zeroInt64, rserrors.ErrInvalidParameter)
			},
			wantErr: amerr.GetErrorsFromCode(amerr.ErrWebServiceTestNotFound).GetErrFromLanguage(langForTest),
		},
		{
			name: "unexpected WebServiceTestService.GetWebServiceTestById error",
			mockFunc: func(mockContext *mocks.Context, mockWebServiceTestService *mocks2.WebServiceTestService, webServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", WebServiceTestIdParam).Return(webServiceTestId, nil)

				mockWebServiceTestService.On("DeleteWebServiceTestById", webServiceTestWithId).
					Return(amerr.GetErrInternalServer())
			},
			wantErr: amerr.GetErrInternalServer().GetErrFromLanguage(langForTest),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContext := &mocks.Context{}
			mockWebServiceTestService := &mocks2.WebServiceTestService{}
			mockWebServiceService := &mocks2.WebServiceService{}

			tt.mockFunc(mockContext, mockWebServiceTestService, mockWebServiceService)
			handler := &WebServiceTestHandlerImpl{
				webServiceService:     mockWebServiceService,
				webServiceTestService: mockWebServiceTestService,
			}
			assert.Equal(t, tt.wantErr, handler.DeleteWebServiceTest(mockContext))
		})
	}
}

func TestNewWebServiceTestHandler(t *testing.T) {
	type args struct {
		webServiceService     services.WebServiceService
		webServiceTestService services.WebServiceTestService
	}
	tests := []struct {
		name    string
		args    args
		want    WebServiceTestHandler
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				webServiceService:     &mocks2.WebServiceService{},
				webServiceTestService: &mocks2.WebServiceTestService{},
			},
			want: &WebServiceTestHandlerImpl{
				webServiceService:     &mocks2.WebServiceService{},
				webServiceTestService: &mocks2.WebServiceTestService{},
			},
			wantErr: false,
		},
		{
			name: "invalid field",
			args: args{
				// webServiceService: &mocks2.WebServiceService{},
				webServiceTestService: &mocks2.WebServiceTestService{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid field",
			args: args{
				webServiceService: &mocks2.WebServiceService{},
				// webServiceTestService:   &mocks2.WebServiceTestService{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewWebServiceTestHandler(tt.args.webServiceService, tt.args.webServiceTestService)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWebServiceTestHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWebServiceTestHandler() got = %v, want %v", got, tt.want)
			}
		})
	}
}
