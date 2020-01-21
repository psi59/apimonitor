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

const (
	TestIdParam = "testId"
)

var _ TestHandler = &TestHandlerImpl{}

type TestHandler interface {
	CreateTest(c echo.Context) error
	GetTest(c echo.Context) error
	DeleteTest(c echo.Context) error
	GetTestList(c echo.Context) error
	UpdateTest(c echo.Context) error
}

type TestHandlerImpl struct {
	webServiceService services.WebServiceService
	testService       services.TestService
}

func (handler *TestHandlerImpl) CreateTest(c echo.Context) error {
	ctx, err := middlewares.ConvertToCustomContext(c)
	if err != nil {
		return errors.WithStack(err)
	}

	lang := ctx.Language()

	webServiceId := ctx.QueryParam(WebServiceIdParam)
	webService := &models.WebService{Id: webServiceId}
	if err := handler.webServiceService.GetWebServiceById(webService); err != nil {
		return err.GetErrFromLanguage(lang)
	}

	var request models.TestRequest
	if err := ctx.Bind(&request); err != nil {
		rslog.Error(err)
		return amerr.GetErrorsFromCode(amerr.ErrBadRequest).GetErrFromLanguage(lang)
	}

	test, aerr := handler.testService.CreateTest(webService, request)
	if aerr != nil {
		return aerr.GetErrFromLanguage(lang)
	}

	return ctx.JSON(http.StatusOK, test)
}

func (handler *TestHandlerImpl) GetTest(c echo.Context) error {
	ctx, err := middlewares.ConvertToCustomContext(c)
	if err != nil {
		return errors.WithStack(err)
	}

	lang := ctx.Language()
	testId := ctx.Param(TestIdParam)
	test := &models.Test{Id: testId}
	if err := handler.testService.GetTestById(test); err != nil {
		return err.GetErrFromLanguage(lang)
	}

	return ctx.JSON(http.StatusOK, test)
}

func (handler *TestHandlerImpl) DeleteTest(c echo.Context) error {
	ctx, err := middlewares.ConvertToCustomContext(c)
	if err != nil {
		return errors.WithStack(err)
	}

	lang := ctx.Language()
	testId := ctx.Param(TestIdParam)
	if rsvalid.IsZero(testId) {
		return amerr.GetErrorsFromCode(amerr.ErrTestNotFound).GetErrFromLanguage(lang)
	}

	test := &models.Test{Id: testId}
	if err := handler.testService.DeleteTestById(test); err != nil {
		return err.GetErrFromLanguage(lang)
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (handler *TestHandlerImpl) GetTestList(c echo.Context) error {
	ctx, err := middlewares.ConvertToCustomContext(c)
	if err != nil {
		return errors.WithStack(err)
	}

	lang := ctx.Language()

	pageInt64, err := ctx.QueryParamInt64("page", 1)
	if err != nil {
		rslog.Error(err)
		return amerr.GetErrorsFromCode(amerr.ErrBadRequest).GetErrFromLanguage(lang)
	}

	numItemInt64, err := ctx.QueryParamInt64("num_item", 20)
	if err != nil {
		rslog.Error(err)
		return amerr.GetErrorsFromCode(amerr.ErrBadRequest).GetErrFromLanguage(lang)
	}

	webServiceId := ctx.Param(WebServiceIdParam)
	list, aerr := handler.testService.GetTestList(models.TestListRequest{
		Page:         int(pageInt64),
		NumItem:      int(numItemInt64),
		WebServiceId: webServiceId,
	})
	if aerr != nil {
		return aerr.GetErrFromLanguage(lang)
	}

	return ctx.JSON(http.StatusOK, list)
}

func (handler *TestHandlerImpl) UpdateTest(c echo.Context) error {
	ctx, err := middlewares.ConvertToCustomContext(c)
	if err != nil {
		return errors.WithStack(err)
	}

	lang := ctx.Language()

	test := &models.Test{
		Id: ctx.Param(TestIdParam),
	}

	var request models.TestRequest
	if err := ctx.Bind(&request); err != nil {
		rslog.Error(err)
		return amerr.GetErrorsFromCode(amerr.ErrBadRequest).GetErrFromLanguage(lang)
	}

	if err := request.Validate(); err != nil {
		rslog.Error(err)
		return amerr.GetErrorsFromCode(amerr.ErrBadRequest).GetErrFromLanguage(lang)
	}

	if err := handler.testService.UpdateTestById(test, request); err != nil {
		return err.GetErrFromLanguage(lang)
	}

	return ctx.JSON(http.StatusOK, nil)
}

func NewTestHandler(webServiceService services.WebServiceService, testService services.TestService) (TestHandler, error) {
	if rsvalid.IsZero(
		webServiceService,
		testService,
	) {
		return nil, errors.Wrap(rserrors.ErrInvalidParameter, "TestHandler")
	}
	return &TestHandlerImpl{
		webServiceService: webServiceService,
		testService:       testService,
	}, nil
}
