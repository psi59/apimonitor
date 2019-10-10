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
	WebServiceTestIdParam = "webServiceTest_id"
)

type WebServiceTestHandler interface {
	CreateWebServiceTest(c echo.Context) error
	GetWebServiceTest(c echo.Context) error
	DeleteWebServiceTest(c echo.Context) error
	GetWebServiceTestList(c echo.Context) error
}

type WebServiceTestHandlerImpl struct {
	webServiceService     services.WebServiceService
	webServiceTestService services.WebServiceTestService
}

func (handler *WebServiceTestHandlerImpl) CreateWebServiceTest(c echo.Context) error {
	ctx, err := middlewares.ConvertToCustomContext(c)
	if err != nil {
		return errors.WithStack(err)
	}

	lang := ctx.Language()
	webServiceId, _ := ctx.ParamInt64(WebServiceIdParam)
	if rsvalid.IsZero(webServiceId) {
		return amerr.GetErrorsFromCode(amerr.ErrWebServiceNotFound).GetErrFromLanguage(lang)
	}

	webService := &models.WebService{Id: webServiceId}
	if err := handler.webServiceService.GetWebServiceById(webService); err != nil {
		return err.GetErrFromLanguage(lang)
	}

	var request models.WebServiceTestRequest
	if err := ctx.Bind(&request); err != nil {
		rslog.Error(err)
		return amerr.GetErrorsFromCode(amerr.ErrBadRequest).GetErrFromLanguage(lang)
	}

	webServiceTest, aerr := handler.webServiceTestService.CreateWebServiceTest(webService, request)
	if aerr != nil {
		return aerr.GetErrFromLanguage(lang)
	}

	return ctx.JSON(http.StatusOK, webServiceTest)
}

func (handler *WebServiceTestHandlerImpl) GetWebServiceTest(c echo.Context) error {
	ctx, err := middlewares.ConvertToCustomContext(c)
	if err != nil {
		return errors.WithStack(err)
	}

	lang := ctx.Language()
	webServiceTestId, _ := ctx.ParamInt64(WebServiceTestIdParam)
	if rsvalid.IsZero(webServiceTestId) {
		return amerr.GetErrorsFromCode(amerr.ErrWebServiceTestNotFound).GetErrFromLanguage(lang)
	}

	webServiceTest := &models.WebServiceTest{Id: webServiceTestId}
	if err := handler.webServiceTestService.GetWebServiceTestById(webServiceTest); err != nil {
		return err.GetErrFromLanguage(lang)
	}

	return ctx.JSON(http.StatusOK, webServiceTest)
}

func (handler *WebServiceTestHandlerImpl) DeleteWebServiceTest(c echo.Context) error {
	ctx, err := middlewares.ConvertToCustomContext(c)
	if err != nil {
		return errors.WithStack(err)
	}

	lang := ctx.Language()
	webServiceTestId, _ := ctx.ParamInt64(WebServiceTestIdParam)
	if rsvalid.IsZero(webServiceTestId) {
		return amerr.GetErrorsFromCode(amerr.ErrWebServiceTestNotFound).GetErrFromLanguage(lang)
	}

	webServiceTest := &models.WebServiceTest{Id: webServiceTestId}
	if err := handler.webServiceTestService.DeleteWebServiceTestById(webServiceTest); err != nil {
		return err.GetErrFromLanguage(lang)
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (handler *WebServiceTestHandlerImpl) GetWebServiceTestList(c echo.Context) error {
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

	webServiceIdInt64, _ := ctx.QueryParamInt64("web_service_id", 0)
	if webServiceIdInt64 == 0 {
		rslog.Error(err)
		return amerr.GetErrorsFromCode(amerr.ErrBadRequest).GetErrFromLanguage(lang)
	}

	list, aerr := handler.webServiceTestService.GetWebServiceTestList(models.WebServiceTestListRequest{
		Page:         int(pageInt64),
		NumItem:      int(numItemInt64),
		WebServiceId: webServiceIdInt64,
	})
	if aerr != nil {
		return aerr.GetErrFromLanguage(lang)
	}

	return ctx.JSON(http.StatusOK, list)
}

func NewWebServiceTestHandler(webServiceService services.WebServiceService, webServiceTestService services.WebServiceTestService) (WebServiceTestHandler, error) {
	if rsvalid.IsZero(
		webServiceService,
		webServiceTestService,
	) {
		return nil, errors.Wrap(rserrors.ErrInvalidParameter, "WebServiceTestHandler")
	}
	return &WebServiceTestHandlerImpl{
		webServiceService:     webServiceService,
		webServiceTestService: webServiceTestService,
	}, nil
}
