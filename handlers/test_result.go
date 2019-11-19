package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/realsangil/apimonitor/middlewares"
	"github.com/realsangil/apimonitor/models"
	"github.com/realsangil/apimonitor/pkg/amerr"
	"github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/rslog"
	"github.com/realsangil/apimonitor/pkg/rsvalid"
	"github.com/realsangil/apimonitor/services"
)

type TestResultHandler interface {
	GetList(c echo.Context) error
}

type TestResultHandlerImpl struct {
	TestResultService services.TestResultService
}

func (handler *TestResultHandlerImpl) GetList(c echo.Context) error {
	ctx, err := middlewares.ConvertToCustomContext(c)
	if err != nil {
		return errors.WithStack(err)
	}

	lang := ctx.Language()
	page, err := ctx.QueryParamInt64("page", 1)
	if err != nil {
		rslog.Error(err)
		return amerr.GetErrorsFromCode(amerr.ErrBadRequest).GetErrFromLanguage(lang)
	}
	numItem, err := ctx.QueryParamInt64("num_item", 20)
	if err != nil {
		rslog.Error(err)
		return amerr.GetErrorsFromCode(amerr.ErrBadRequest).GetErrFromLanguage(lang)
	}
	isSuccess := ctx.QueryParam("is_success")

	testId, err := ctx.ParamInt64(TestIdParam)
	if err != nil {
		rslog.Error(err)
		return amerr.GetErrorsFromCode(amerr.ErrWebServiceTestNotFound).GetErrFromLanguage(lang)
	}

	request := models.TestResultListRequest{
		Page:      int(page),
		NumItem:   int(numItem),
		IsSuccess: models.IsSuccess(isSuccess),
		// StartTestedAt: time.Time{},
		// EndTestedAt:   time.Time{},
	}
	test := &models.WebServiceTest{Id: testId}

	list, aerr := handler.TestResultService.GetResultListByTestId(test, request)
	if aerr != nil {
		return aerr.GetErrFromLanguage(lang)
	}

	return ctx.JSON(http.StatusOK, list)
}

func NewTestResultHandler(testResultService services.TestResultService) (TestResultHandler, error) {
	if rsvalid.IsZero(testResultService) {
		return nil, rserrors.ErrInvalidParameter
	}
	return &TestResultHandlerImpl{TestResultService: testResultService}, nil
}
