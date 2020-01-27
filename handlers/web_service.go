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
	WebServiceIdParam = "webServiceId"
)

type WebServiceHandler interface {
	CreateWebService(c echo.Context) error
	GetWebServiceById(c echo.Context) error
	DeleteWebServiceById(c echo.Context) error
	UpdateWebServiceById(c echo.Context) error
	GetWebServiceList(c echo.Context) error
	ExecuteTests(c echo.Context) error
}

type WebServiceHandlerImpl struct {
	webServiceService services.WebServiceService
}

func (handler *WebServiceHandlerImpl) CreateWebService(c echo.Context) error {
	ctx, err := middlewares.ConvertToCustomContext(c)
	if err != nil {
		return errors.WithStack(err)
	}

	lang := ctx.Language()

	var request models.WebServiceRequest
	if err := ctx.Bind(&request); err != nil {
		rslog.Error(err)
		return amerr.GetErrorsFromCode(amerr.ErrBadRequest).GetErrFromLanguage(lang)
	}
	webService, aerr := handler.webServiceService.CreateWebService(request)
	if aerr != nil {
		return aerr.GetErrFromLanguage(lang)
	}

	return ctx.JSON(http.StatusOK, webService)
}

func (handler *WebServiceHandlerImpl) GetWebServiceById(c echo.Context) error {
	ctx, err := middlewares.ConvertToCustomContext(c)
	if err != nil {
		return errors.WithStack(err)
	}

	lang := ctx.Language()

	webServiceId := ctx.Param(WebServiceIdParam)
	webService := &models.WebService{Id: webServiceId}
	if err := handler.webServiceService.GetWebServiceById(webService); err != nil {
		return err.GetErrFromLanguage(lang)
	}

	return ctx.JSON(http.StatusOK, webService)
}

func (handler *WebServiceHandlerImpl) DeleteWebServiceById(c echo.Context) error {
	ctx, err := middlewares.ConvertToCustomContext(c)
	if err != nil {
		return errors.WithStack(err)
	}

	lang := ctx.Language()

	webServiceId := ctx.Param(WebServiceIdParam)
	if err != nil {
		rslog.Error(err)
		return amerr.GetErrorsFromCode(amerr.ErrWebServiceNotFound).GetErrFromLanguage(lang)
	}

	webService := &models.WebService{Id: webServiceId}
	if err := handler.webServiceService.DeleteWebServiceById(webService); err != nil {
		return err.GetErrFromLanguage(lang)
	}

	return ctx.JSON(http.StatusOK, webService)
}

func (handler *WebServiceHandlerImpl) UpdateWebServiceById(c echo.Context) error {
	ctx, err := middlewares.ConvertToCustomContext(c)
	if err != nil {
		return errors.WithStack(err)
	}

	lang := ctx.Language()

	webServiceId := ctx.Param(WebServiceIdParam)
	webService := &models.WebService{Id: webServiceId}
	var request models.WebServiceRequest
	if err := ctx.Bind(&request); err != nil {
		rslog.Error(err)
		return amerr.GetErrorsFromCode(amerr.ErrBadRequest).GetErrFromLanguage(lang)
	}

	if err := handler.webServiceService.UpdateWebServiceById(webService, request); err != nil {
		return err.GetErrFromLanguage(lang)
	}

	return ctx.JSON(http.StatusOK, webService)
}

func (handler *WebServiceHandlerImpl) GetWebServiceList(c echo.Context) error {
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

	list, aerr := handler.webServiceService.GetWebServiceList(models.WebServiceListRequest{
		Page:    int(page),
		NumItem: int(numItem),
	})
	if aerr != nil {
		return aerr.GetErrFromLanguage(lang)
	}

	return ctx.JSON(http.StatusOK, list)
}

func (handler *WebServiceHandlerImpl) ExecuteTests(c echo.Context) error {
	ctx, err := middlewares.ConvertToCustomContext(c)
	if err != nil {
		return errors.WithStack(err)
	}

	lang := ctx.Language()

	webService := &models.WebService{
		Id: ctx.Param(WebServiceIdParam),
	}

	if err := handler.webServiceService.ExecuteTests(webService); err != nil {
		rslog.Error(err)
		return err.GetErrFromLanguage(lang)
	}

	return ctx.JSON(http.StatusOK, nil)
}

func NewWebServiceHandler(webServiceService services.WebServiceService) (WebServiceHandler, error) {
	if rsvalid.IsZero(webServiceService) {
		return nil, rserrors.ErrInvalidParameter
	}
	return &WebServiceHandlerImpl{
		webServiceService: webServiceService,
	}, nil
}
