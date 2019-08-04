package services

import (
	"reflect"
	"testing"
	"time"

	"github.com/realsangil/apimonitor/models"
	"github.com/realsangil/apimonitor/pkg/amerr"
	"github.com/realsangil/apimonitor/pkg/rsdb"
	mocks2 "github.com/realsangil/apimonitor/pkg/rsdb/mocks"
	"github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/testutils"
	"github.com/realsangil/apimonitor/repositories"
	"github.com/realsangil/apimonitor/repositories/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type webServiceMockFunc func(mockWebServiceRepository *mocks.WebServiceRepository)

func TestWebServiceServiceImpl_CreateWebService(t *testing.T) {
	testutils.MonkeyAll()

	mockTx := &mocks2.Transaction{}
	validatedRequest := models.WebServiceRequest{
		Host:    "https://realsangil.github.io",
		Desc:    "sangil's dev blog",
		Favicon: "",
	}

	validatedWebServiceWithoutId := &models.WebService{
		DefaultValidateChecker: models.ValidatedDefaultValidateChecker,
		Host:                   "realsangil.github.io",
		HttpSchema:             "https",
		Desc:                   "sangil's dev blog",
		Favicon:                "",
		Created:                time.Now(),
		LastModified:           time.Now(),
	}

	validatedWebService := &models.WebService{
		DefaultValidateChecker: models.ValidatedDefaultValidateChecker,
		Id:                     1,
		Host:                   "realsangil.github.io",
		HttpSchema:             "https",
		Desc:                   "sangil's dev blog",
		Favicon:                "",
		Created:                time.Now(),
		LastModified:           time.Now(),
	}

	mockFuncWithError := func(err error) webServiceMockFunc {
		return func(mockWebServiceRepository *mocks.WebServiceRepository) {
			m := mockWebServiceRepository.On("Create", mockTx, validatedWebServiceWithoutId)
			if err == nil {
				m.Run(func(args mock.Arguments) {
					arg := args.Get(1).(*models.WebService)
					arg.Id = 1
				})
			}
			m.Return(err)
		}
	}

	type args struct {
		transaction rsdb.Transaction
		request     models.WebServiceRequest
	}
	tests := []struct {
		name     string
		args     args
		mockFunc webServiceMockFunc
		want     *models.WebService
		wantErr  *amerr.ErrorWithLanguage
	}{
		{
			name: "pass_https_host",
			args: args{
				transaction: mockTx,
				request:     validatedRequest,
			},
			mockFunc: mockFuncWithError(nil),
			want:     validatedWebService,
			wantErr:  nil,
		},
		{
			name: "invalid_host",
			args: args{
				transaction: mockTx,
				request: models.WebServiceRequest{
					Host:    "ftp://realsangil.github.io",
					Desc:    "sangil's dev blog",
					Favicon: "",
				},
			},
			mockFunc: func(mockWebServiceRepository *mocks.WebServiceRepository) {
			},
			want:    nil,
			wantErr: amerr.GetErrorsFromCode(amerr.ErrBadRequest),
		},
		{
			name: "duplicated_web_service",
			args: args{
				transaction: mockTx,
				request:     validatedRequest,
			},
			mockFunc: mockFuncWithError(rsdb.ErrDuplicateData),
			want:     nil,
			wantErr:  amerr.GetErrorsFromCode(amerr.ErrDuplicatedWebService),
		},
		{
			name: "data too long",
			args: args{
				transaction: mockTx,
				request:     validatedRequest,
			},
			mockFunc: mockFuncWithError(rsdb.ErrInvalidData),
			want:     nil,
			wantErr:  amerr.GetErrorsFromCode(amerr.ErrBadRequest),
		},
		{
			name: "unexpected error",
			args: args{
				transaction: mockTx,
				request:     validatedRequest,
			},
			mockFunc: mockFuncWithError(rserrors.ErrUnexpected),
			want:     nil,
			wantErr:  amerr.GetErrorsFromCode(amerr.ErrInternalServer),
		},
		{
			name: "invalid transaction",
			args: args{
				transaction: nil,
				request:     validatedRequest,
			},
			mockFunc: func(mockWebServiceRepository *mocks.WebServiceRepository) {},
			want:     nil,
			wantErr:  amerr.GetErrorsFromCode(amerr.ErrInternalServer),
		},
		{
			name: "invalid transaction",
			args: args{
				transaction: rsdb.Transaction(nil),
				request:     validatedRequest,
			},
			mockFunc: func(mockWebServiceRepository *mocks.WebServiceRepository) {},
			want:     nil,
			wantErr:  amerr.GetErrorsFromCode(amerr.ErrInternalServer),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWebServiceRepository := &mocks.WebServiceRepository{}
			service := &WebServiceServiceImpl{
				webServiceRepository: mockWebServiceRepository,
			}
			tt.mockFunc(mockWebServiceRepository)
			got, gotErr := service.CreateWebService(tt.args.transaction, tt.args.request)
			assert.Equal(t, tt.wantErr, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewWebServiceService(t *testing.T) {
	testutils.MonkeyAll()

	mockWebServiceRepository := &mocks.WebServiceRepository{}

	type args struct {
		webServiceRepository repositories.WebServiceRepository
	}
	tests := []struct {
		name    string
		args    args
		want    WebServiceService
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				webServiceRepository: mockWebServiceRepository,
			},
			want: &WebServiceServiceImpl{
				webServiceRepository: mockWebServiceRepository,
			},
			wantErr: false,
		},
		{
			name: "invalid WebServiceRepository",
			args: args{
				webServiceRepository: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid WebServiceRepository",
			args: args{
				webServiceRepository: repositories.WebServiceRepository(nil),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewWebServiceService(tt.args.webServiceRepository)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWebServiceService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWebServiceService() = %v, want %v", got, tt.want)
			}
		})
	}
}
