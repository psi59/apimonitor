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
	"github.com/realsangil/apimonitor/pkg/rsmodels"
	"github.com/realsangil/apimonitor/pkg/testutils"
	"github.com/realsangil/apimonitor/services"
	mocks2 "github.com/realsangil/apimonitor/services/mocks"
)

const langForTest = "ko"

func setContextLanguageForTest(ctx *mocks.Context) {
	ctx.On("Language").Return(langForTest)
}

type testMockFunc func(mockContext *mocks.Context, mockTestService *mocks2.TestService, mockWebServiceService *mocks2.WebServiceService)

func TestTestHandlerImpl_CreateTest(t *testing.T) {
	testutils.MonkeyAll()

	webServiceId := int64(1)
	zeroInt64 := int64(0)

	webServiceWithId := &models.WebService{Id: webServiceId}
	webService := &models.WebService{
		DefaultValidateChecker: rsmodels.ValidatedDefaultValidateChecker,
		Id:                     1,
		Host:                   "realsangil.github.io",
		HttpSchema:             "https",
		Desc:                   "sangil's dev blog",
		Favicon:                "",
		Created:                time.Now(),
		LastModified:           time.Now(),
	}

	request := models.TestRequest{
		Path:        "/path/to/uri",
		HttpMethod:  rshttp.MethodGet,
		ContentType: rshttp.MIMEApplicationJSON,
		RequestData: nil,
		Header:      nil,
		QueryParam:  nil,
	}

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

	tests := []struct {
		name     string
		mockFunc testMockFunc
		wantErr  error
	}{
		{
			name: "pass",
			mockFunc: func(mockContext *mocks.Context, mockTestService *mocks2.TestService, mockWebServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", WebServiceIdParam).Return(webServiceId, nil)

				mockWebServiceService.On("GetWebServiceById", webServiceWithId).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*models.WebService)
						*arg = *webService
					}).
					Return(nil)

				mockContext.On("Bind", &models.TestRequest{}).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*models.TestRequest)
						*arg = request
					}).
					Return(nil)

				mockTestService.On("CreateTest", webService, request).
					Return(test, nil)

				mockContext.On("JSON", http.StatusOK, test).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "invalid webServiceId",
			mockFunc: func(mockContext *mocks.Context, mockTestService *mocks2.TestService, mockWebServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", WebServiceIdParam).Return(zeroInt64, rserrors.ErrInvalidParameter)
			},
			wantErr: amerr.GetErrorsFromCode(amerr.ErrWebServiceNotFound).GetErrFromLanguage(langForTest),
		},
		{
			name: "unexpected webServiceService.GetWebServiceById error",
			mockFunc: func(mockContext *mocks.Context, mockTestService *mocks2.TestService, mockWebServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", WebServiceIdParam).Return(webServiceId, nil)

				mockWebServiceService.On("GetWebServiceById", webServiceWithId).
					Return(amerr.GetErrInternalServer())
			},
			wantErr: amerr.GetErrInternalServer().GetErrFromLanguage(langForTest),
		},
		{
			name: "bind error",
			mockFunc: func(mockContext *mocks.Context, mockTestService *mocks2.TestService, mockWebServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", WebServiceIdParam).Return(webServiceId, nil)

				mockWebServiceService.On("GetWebServiceById", webServiceWithId).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*models.WebService)
						*arg = *webService
					}).
					Return(nil)

				mockContext.On("Bind", &models.TestRequest{}).
					Return(rserrors.ErrUnexpected)
			},
			wantErr: amerr.GetErrorsFromCode(amerr.ErrBadRequest).GetErrFromLanguage(langForTest),
		},
		{
			name: "unexpected TestService.CreateTest error",
			mockFunc: func(mockContext *mocks.Context, mockTestService *mocks2.TestService, mockWebServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", WebServiceIdParam).Return(webServiceId, nil)

				mockWebServiceService.On("GetWebServiceById", webServiceWithId).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*models.WebService)
						*arg = *webService
					}).
					Return(nil)

				mockContext.On("Bind", &models.TestRequest{}).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*models.TestRequest)
						*arg = request
					}).
					Return(nil)

				mockTestService.On("CreateTest", webService, request).
					Return(nil, amerr.GetErrInternalServer())
			},
			wantErr: amerr.GetErrInternalServer().GetErrFromLanguage(langForTest),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContext := &mocks.Context{}
			mockTestService := &mocks2.TestService{}
			mockWebServiceService := &mocks2.WebServiceService{}

			tt.mockFunc(mockContext, mockTestService, mockWebServiceService)
			handler := &TestHandlerImpl{
				webServiceService: mockWebServiceService,
				testService:       mockTestService,
			}
			assert.Equal(t, tt.wantErr, handler.CreateTest(mockContext))
		})
	}
}

func TestTestHandlerImpl_GetTest(t *testing.T) {
	testutils.MonkeyAll()

	zeroInt64 := int64(0)
	testId := int64(1)

	testWithId := &models.Test{Id: testId}
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

	tests := []struct {
		name     string
		mockFunc testMockFunc
		wantErr  error
	}{
		{
			name: "pass",
			mockFunc: func(mockContext *mocks.Context, mockTestService *mocks2.TestService, mockWebServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", TestIdParam).Return(testId, nil)

				mockTestService.On("GetTestById", testWithId).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*models.Test)
						*arg = *test
					}).
					Return(nil)

				mockContext.On("JSON", http.StatusOK, test).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "invalid test id parameter",
			mockFunc: func(mockContext *mocks.Context, mockTestService *mocks2.TestService, webServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", TestIdParam).Return(zeroInt64, rserrors.ErrInvalidParameter)
			},
			wantErr: amerr.GetErrorsFromCode(amerr.ErrTestNotFound).GetErrFromLanguage(langForTest),
		},
		{
			name: "unexpected TestService.GetTestById error",
			mockFunc: func(mockContext *mocks.Context, mockTestService *mocks2.TestService, webServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", TestIdParam).Return(testId, nil)

				mockTestService.On("GetTestById", testWithId).
					Return(amerr.GetErrInternalServer())
			},
			wantErr: amerr.GetErrInternalServer().GetErrFromLanguage(langForTest),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContext := &mocks.Context{}
			mockTestService := &mocks2.TestService{}
			mockWebServiceService := &mocks2.WebServiceService{}

			tt.mockFunc(mockContext, mockTestService, mockWebServiceService)
			handler := &TestHandlerImpl{
				webServiceService: mockWebServiceService,
				testService:       mockTestService,
			}
			assert.Equal(t, tt.wantErr, handler.GetTest(mockContext))
		})
	}
}

func TestTestHandlerImpl_DeleteTest(t *testing.T) {
	testutils.MonkeyAll()

	zeroInt64 := int64(0)
	testId := int64(1)

	testWithId := &models.Test{Id: testId}

	tests := []struct {
		name     string
		mockFunc testMockFunc
		wantErr  error
	}{
		{
			name: "pass",
			mockFunc: func(mockContext *mocks.Context, mockTestService *mocks2.TestService, mockWebServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", TestIdParam).Return(testId, nil)

				mockTestService.On("DeleteTestById", testWithId).
					Return(nil)

				mockContext.On("JSON", http.StatusOK, nil).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "invalid test id parameter",
			mockFunc: func(mockContext *mocks.Context, mockTestService *mocks2.TestService, webServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", TestIdParam).Return(zeroInt64, rserrors.ErrInvalidParameter)
			},
			wantErr: amerr.GetErrorsFromCode(amerr.ErrTestNotFound).GetErrFromLanguage(langForTest),
		},
		{
			name: "unexpected TestService.GetTestById error",
			mockFunc: func(mockContext *mocks.Context, mockTestService *mocks2.TestService, webServiceService *mocks2.WebServiceService) {
				setContextLanguageForTest(mockContext)
				mockContext.On("ParamInt64", TestIdParam).Return(testId, nil)

				mockTestService.On("DeleteTestById", testWithId).
					Return(amerr.GetErrInternalServer())
			},
			wantErr: amerr.GetErrInternalServer().GetErrFromLanguage(langForTest),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContext := &mocks.Context{}
			mockTestService := &mocks2.TestService{}
			mockWebServiceService := &mocks2.WebServiceService{}

			tt.mockFunc(mockContext, mockTestService, mockWebServiceService)
			handler := &TestHandlerImpl{
				webServiceService: mockWebServiceService,
				testService:       mockTestService,
			}
			assert.Equal(t, tt.wantErr, handler.DeleteTest(mockContext))
		})
	}
}

func TestNewTestHandler(t *testing.T) {
	type args struct {
		webServiceService services.WebServiceService
		testService       services.TestService
	}
	tests := []struct {
		name    string
		args    args
		want    TestHandler
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				webServiceService: &mocks2.WebServiceService{},
				testService:       &mocks2.TestService{},
			},
			want: &TestHandlerImpl{
				webServiceService: &mocks2.WebServiceService{},
				testService:       &mocks2.TestService{},
			},
			wantErr: false,
		},
		{
			name: "invalid field",
			args: args{
				// webServiceService: &mocks2.WebServiceService{},
				testService: &mocks2.TestService{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid field",
			args: args{
				webServiceService: &mocks2.WebServiceService{},
				// testService:   &mocks2.TestService{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTestHandler(tt.args.webServiceService, tt.args.testService)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTestHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTestHandler() got = %v, want %v", got, tt.want)
			}
		})
	}
}
