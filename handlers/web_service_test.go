package handlers

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/lalaworks/kmessenger-webserver-go/pkg/lalaerrors"
	mocks2 "github.com/realsangil/apimonitor/middlewares/mocks"
	"github.com/realsangil/apimonitor/models"
	"github.com/realsangil/apimonitor/pkg/amerr"
	mocks3 "github.com/realsangil/apimonitor/pkg/rsdb/mocks"
	"github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/testutils"
	"github.com/realsangil/apimonitor/services"
	"github.com/realsangil/apimonitor/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var lang = "ko"

type webServiceHandlerMockFunc func(mockContext *mocks2.Context, mockWebServiceService *mocks.WebServiceService)

func TestWebServiceHandlerImpl_CreateWebService(t *testing.T) {
	testutils.MonkeyAll()

	mockTx := &mocks3.Transaction{}

	tests := []struct {
		name     string
		mockFunc webServiceHandlerMockFunc
		wantErr  error
	}{
		{
			name: "pass",
			mockFunc: func(mockContext *mocks2.Context, mockWebServiceService *mocks.WebServiceService) {
				request := models.WebServiceRequest{
					Host:    "https://realsangil.github.io",
					Desc:    "sangil's blog",
					Favicon: "",
				}
				webService := &models.WebService{
					DefaultValidateChecker: models.ValidatedDefaultValidateChecker,
					Id:                     1,
					Host:                   "realsangil.github.io",
					HttpSchema:             "https",
					Desc:                   "sangil's blog",
					Favicon:                "",
					Created:                time.Now(),
					LastModified:           time.Now(),
				}

				mockContext.On("Language").Return(lang)
				mockContext.On("GetTx").Return(mockTx, nil)
				mockContext.On("Bind", &models.WebServiceRequest{}).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*models.WebServiceRequest)
					*arg = request
				}).Return(nil)

				mockWebServiceService.On("CreateWebService", mockTx, request).Return(webService, nil)

				mockContext.On("JSON", http.StatusOK, webService).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "transaction error",
			mockFunc: func(mockContext *mocks2.Context, mockWebServiceService *mocks.WebServiceService) {
				mockContext.On("Language").Return(lang)
				mockContext.On("GetTx").Return(nil, rserrors.ErrUnexpected)
			},
			wantErr: amerr.GetErrInternalServer().GetErrFromLanguage(lang),
		},
		{
			name: "bind error",
			mockFunc: func(mockContext *mocks2.Context, mockWebServiceService *mocks.WebServiceService) {
				mockContext.On("Language").Return(lang)
				mockContext.On("GetTx").Return(mockTx, nil)
				mockContext.On("Bind", &models.WebServiceRequest{}).Return(rserrors.ErrUnexpected)
			},
			wantErr: amerr.GetErrorsFromCode(amerr.ErrBadRequest).GetErrFromLanguage(lang),
		},
		{
			name: "service error",
			mockFunc: func(mockContext *mocks2.Context, mockWebServiceService *mocks.WebServiceService) {
				request := models.WebServiceRequest{
					Host:    "https://realsangil.github.io",
					Desc:    "sangil's blog",
					Favicon: "",
				}

				mockContext.On("Language").Return("ko")
				mockContext.On("GetTx").Return(mockTx, nil)
				mockContext.On("Bind", &models.WebServiceRequest{}).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*models.WebServiceRequest)
					*arg = request
				}).Return(nil)

				mockWebServiceService.On("CreateWebService", mockTx, request).Return(nil, amerr.GetErrInternalServer())
			},
			wantErr: amerr.GetErrInternalServer().GetErrFromLanguage(lang),
		},
		{
			name: "service error",
			mockFunc: func(mockContext *mocks2.Context, mockWebServiceService *mocks.WebServiceService) {
				request := models.WebServiceRequest{
					Host:    "https://realsangil.github.io",
					Desc:    "sangil's blog",
					Favicon: "",
				}

				mockContext.On("Language").Return("ko")
				mockContext.On("GetTx").Return(mockTx, nil)
				mockContext.On("Bind", &models.WebServiceRequest{}).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*models.WebServiceRequest)
					*arg = request
				}).Return(nil)

				mockWebServiceService.On("CreateWebService", mockTx, request).Return(nil, amerr.GetErrorsFromCode(amerr.ErrDuplicatedWebService))
			},
			wantErr: amerr.GetErrorsFromCode(amerr.ErrDuplicatedWebService).GetErrFromLanguage(lang),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWebServiceService := &mocks.WebServiceService{}
			mockContext := &mocks2.Context{}

			handler := &WebServiceHandlerImpl{
				webServiceService: mockWebServiceService,
			}

			tt.mockFunc(mockContext, mockWebServiceService)

			gotErr := handler.CreateWebService(mockContext)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}

func TestNewWebServiceHandler(t *testing.T) {
	testutils.MonkeyAll()

	mockWebServiceService := &mocks.WebServiceService{}

	type args struct {
		webServiceService services.WebServiceService
	}
	tests := []struct {
		name    string
		args    args
		want    WebServiceHandler
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				webServiceService: mockWebServiceService,
			},
			want: &WebServiceHandlerImpl{
				webServiceService: mockWebServiceService,
			},
			wantErr: false,
		},
		{
			name: "invalid service",
			args: args{
				webServiceService: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid service",
			args: args{
				webServiceService: services.WebServiceService(nil),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewWebServiceHandler(tt.args.webServiceService)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWebServiceHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWebServiceHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebServiceHandlerImpl_GetWebServiceById(t *testing.T) {
	testutils.MonkeyAll()

	mockTx := &mocks3.Transaction{}

	tests := []struct {
		name     string
		mockFunc webServiceHandlerMockFunc
		wantErr  error
	}{
		{
			name: "pass",
			mockFunc: func(mockContext *mocks2.Context, mockWebServiceService *mocks.WebServiceService) {
				webService := &models.WebService{
					Id:           1,
					Host:         "realsangil.github.io",
					HttpSchema:   "https",
					Desc:         "sangil's dev blog",
					Favicon:      "",
					Created:      time.Now(),
					LastModified: time.Now(),
				}
				mockContext.On("Language").Return(lang)
				mockContext.On("GetTx").Return(mockTx, nil)
				mockContext.On("ParamInt64", WebServiceIdParam).Return(int64(1), nil)

				mockWebServiceService.On("GetWebServiceById", mockTx, &models.WebService{Id: 1}).Run(func(args mock.Arguments) {
					arg := args.Get(1).(*models.WebService)
					*arg = *webService
				}).Return(nil)

				mockContext.On("JSON", http.StatusOK, webService).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "transaction error",
			mockFunc: func(mockContext *mocks2.Context, mockWebServiceService *mocks.WebServiceService) {
				mockContext.On("Language").Return(lang)
				mockContext.On("GetTx").Return(nil, rserrors.ErrUnexpected)
				mockContext.On("ParamInt64", WebServiceIdParam).Return(int64(1), nil)
			},
			wantErr: amerr.GetErrInternalServer().GetErrFromLanguage(lang),
		},
		{
			name: "invalid webServiceId",
			mockFunc: func(mockContext *mocks2.Context, mockWebServiceService *mocks.WebServiceService) {
				mockContext.On("Language").Return(lang)
				mockContext.On("GetTx").Return(mockTx, nil)
				mockContext.On("ParamInt64", WebServiceIdParam).Return(int64(0), rserrors.ErrUnexpected)
			},
			wantErr: amerr.GetErrorsFromCode(amerr.ErrWebServiceNotFound).GetErrFromLanguage(lang),
		},
		{
			name: "unexpected Service error",
			mockFunc: func(mockContext *mocks2.Context, mockWebServiceService *mocks.WebServiceService) {
				mockContext.On("Language").Return(lang)
				mockContext.On("GetTx").Return(mockTx, nil)
				mockContext.On("ParamInt64", "web_service_id").Return(int64(1), nil)

				mockWebServiceService.On("GetWebServiceById", mockTx, &models.WebService{Id: 1}).Return(amerr.GetErrInternalServer())
			},
			wantErr: amerr.GetErrInternalServer().GetErrFromLanguage(lang),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContext := &mocks2.Context{}
			mockWebServiceService := &mocks.WebServiceService{}
			handler := &WebServiceHandlerImpl{
				webServiceService: mockWebServiceService,
			}
			tt.mockFunc(mockContext, mockWebServiceService)
			err := handler.GetWebServiceById(mockContext)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestWebServiceHandlerImpl_DeleteWebServiceById(t *testing.T) {
	testutils.MonkeyAll()

	mockTx := &mocks3.Transaction{}

	tests := []struct {
		name     string
		mockFunc webServiceHandlerMockFunc
		wantErr  error
	}{
		{
			name: "pass",
			mockFunc: func(mockContext *mocks2.Context, mockWebServiceService *mocks.WebServiceService) {
				webService := &models.WebService{
					Id:           1,
					Host:         "realsangil.github.io",
					HttpSchema:   "https",
					Desc:         "sangil's dev blog",
					Favicon:      "",
					Created:      time.Now(),
					LastModified: time.Now(),
				}
				mockContext.On("Language").Return(lang)
				mockContext.On("GetTx").Return(mockTx, nil)
				mockContext.On("ParamInt64", WebServiceIdParam).Return(int64(1), nil)

				mockWebServiceService.On("DeleteWebServiceById", mockTx, &models.WebService{Id: 1}).Run(func(args mock.Arguments) {
					arg := args.Get(1).(*models.WebService)
					*arg = *webService
				}).Return(nil)

				mockContext.On("JSON", http.StatusOK, webService).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "transaction error",
			mockFunc: func(mockContext *mocks2.Context, mockWebServiceService *mocks.WebServiceService) {
				mockContext.On("Language").Return(lang)
				mockContext.On("GetTx").Return(nil, rserrors.ErrUnexpected)
				mockContext.On("ParamInt64", WebServiceIdParam).Return(int64(1), nil)
			},
			wantErr: amerr.GetErrInternalServer().GetErrFromLanguage(lang),
		},
		{
			name: "invalid webServiceId",
			mockFunc: func(mockContext *mocks2.Context, mockWebServiceService *mocks.WebServiceService) {
				mockContext.On("Language").Return(lang)
				mockContext.On("GetTx").Return(mockTx, nil)
				mockContext.On("ParamInt64", WebServiceIdParam).Return(int64(0), rserrors.ErrUnexpected)
			},
			wantErr: amerr.GetErrorsFromCode(amerr.ErrWebServiceNotFound).GetErrFromLanguage(lang),
		},
		{
			name: "unexpected Service error",
			mockFunc: func(mockContext *mocks2.Context, mockWebServiceService *mocks.WebServiceService) {
				mockContext.On("Language").Return(lang)
				mockContext.On("GetTx").Return(mockTx, nil)
				mockContext.On("ParamInt64", "web_service_id").Return(int64(1), nil)

				mockWebServiceService.On("DeleteWebServiceById", mockTx, &models.WebService{Id: 1}).Return(amerr.GetErrInternalServer())
			},
			wantErr: amerr.GetErrInternalServer().GetErrFromLanguage(lang),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContext := &mocks2.Context{}
			mockWebServiceService := &mocks.WebServiceService{}

			handler := &WebServiceHandlerImpl{
				webServiceService: mockWebServiceService,
			}

			tt.mockFunc(mockContext, mockWebServiceService)
			err := handler.DeleteWebServiceById(mockContext)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestWebServiceHandlerImpl_UpdateWebServiceById(t *testing.T) {
	testutils.MonkeyAll()

	mockTx := &mocks3.Transaction{}

	webServiceWithId := &models.WebService{Id: 1}
	webService := &models.WebService{
		Id:           1,
		Host:         "realsangil.github.io",
		HttpSchema:   "https",
		Desc:         "sangil's dev blog",
		Favicon:      "",
		Created:      time.Now(),
		LastModified: time.Now(),
	}
	webServiceRequest := models.WebServiceRequest{
		Host:    "https://www.lalaworks.com",
		Desc:    "lalaworks website",
		Favicon: "",
	}

	tests := []struct {
		name     string
		mockFunc webServiceHandlerMockFunc
		wantErr  error
	}{
		{
			name: "pass",
			mockFunc: func(mockContext *mocks2.Context, mockWebServiceService *mocks.WebServiceService) {
				mockContext.On("Language").Return(lang)
				mockContext.On("GetTx").Return(mockTx, nil)
				mockContext.On("ParamInt64", WebServiceIdParam).Return(int64(1), nil)
				mockContext.On("Bind", &models.WebServiceRequest{}).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*models.WebServiceRequest)
					*arg = webServiceRequest
				}).Return(nil)

				mockWebServiceService.On("UpdateWebServiceById", mockTx, webServiceWithId, webServiceRequest).Run(func(args mock.Arguments) {
					arg := args.Get(1).(*models.WebService)
					*arg = *webService
				}).Return(nil)

				mockContext.On("JSON", http.StatusOK, webService).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "transaction error",
			mockFunc: func(mockContext *mocks2.Context, mockWebServiceService *mocks.WebServiceService) {
				mockContext.On("Language").Return(lang)
				mockContext.On("GetTx").Return(nil, rserrors.ErrUnexpected)
			},
			wantErr: amerr.GetErrInternalServer().GetErrFromLanguage(lang),
		},
		{
			name: "invalid WebServiceId",
			mockFunc: func(mockContext *mocks2.Context, mockWebServiceService *mocks.WebServiceService) {
				mockContext.On("Language").Return(lang)
				mockContext.On("GetTx").Return(mockTx, nil)
				mockContext.On("ParamInt64", WebServiceIdParam).Return(int64(0), rserrors.ErrUnexpected)
			},
			wantErr: amerr.GetErrorsFromCode(amerr.ErrWebServiceNotFound).GetErrFromLanguage(lang),
		},
		{
			name: "Bind error",
			mockFunc: func(mockContext *mocks2.Context, mockWebServiceService *mocks.WebServiceService) {
				mockContext.On("Language").Return(lang)
				mockContext.On("GetTx").Return(mockTx, nil)
				mockContext.On("ParamInt64", WebServiceIdParam).Return(int64(1), nil)
				mockContext.On("Bind", &models.WebServiceRequest{}).Return(lalaerrors.ErrUnexpected)
			},
			wantErr: amerr.GetErrorsFromCode(amerr.ErrBadRequest).GetErrFromLanguage(lang),
		},
		{
			name: "WebServiceService.UpdateWebServiceById error",
			mockFunc: func(mockContext *mocks2.Context, mockWebServiceService *mocks.WebServiceService) {
				mockContext.On("Language").Return(lang)
				mockContext.On("GetTx").Return(mockTx, nil)
				mockContext.On("ParamInt64", WebServiceIdParam).Return(int64(1), nil)
				mockContext.On("Bind", &models.WebServiceRequest{}).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*models.WebServiceRequest)
					*arg = webServiceRequest
				}).Return(nil)

				mockWebServiceService.On("UpdateWebServiceById", mockTx, webServiceWithId, webServiceRequest).Return(amerr.GetErrInternalServer())
			},
			wantErr: amerr.GetErrInternalServer().GetErrFromLanguage(lang),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContext := &mocks2.Context{}
			mockWebServiceService := &mocks.WebServiceService{}

			handler := &WebServiceHandlerImpl{
				webServiceService: mockWebServiceService,
			}

			tt.mockFunc(mockContext, mockWebServiceService)
			err := handler.UpdateWebServiceById(mockContext)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
