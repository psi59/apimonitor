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
	EndpointIdParam = "endpoint_id"
)

type EndpointHandler interface {
	CreateEndpoint(c echo.Context) error
	GetEndpoint(c echo.Context) error
	DeleteEndpoint(c echo.Context) error
	GetEndpointList(c echo.Context) error
}

type EndpointHandlerImpl struct {
	webServiceService services.WebServiceService
	endpointService   services.EndpointService
}

func (handler *EndpointHandlerImpl) CreateEndpoint(c echo.Context) error {
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

	var request models.EndpointRequest
	if err := ctx.Bind(&request); err != nil {
		rslog.Error(err)
		return amerr.GetErrorsFromCode(amerr.ErrBadRequest).GetErrFromLanguage(lang)
	}

	endpoint, aerr := handler.endpointService.CreateEndpoint(webService, request)
	if aerr != nil {
		return aerr.GetErrFromLanguage(lang)
	}

	return ctx.JSON(http.StatusOK, endpoint)
}

func (handler *EndpointHandlerImpl) GetEndpoint(c echo.Context) error {
	ctx, err := middlewares.ConvertToCustomContext(c)
	if err != nil {
		return errors.WithStack(err)
	}

	lang := ctx.Language()
	endpointId, _ := ctx.ParamInt64(EndpointIdParam)
	if rsvalid.IsZero(endpointId) {
		return amerr.GetErrorsFromCode(amerr.ErrEndpointNotFound).GetErrFromLanguage(lang)
	}

	endpoint := &models.Endpoint{Id: endpointId}
	if err := handler.endpointService.GetEndpointById(endpoint); err != nil {
		return err.GetErrFromLanguage(lang)
	}

	return ctx.JSON(http.StatusOK, endpoint)
}

func (handler *EndpointHandlerImpl) DeleteEndpoint(c echo.Context) error {
	ctx, err := middlewares.ConvertToCustomContext(c)
	if err != nil {
		return errors.WithStack(err)
	}

	lang := ctx.Language()
	endpointId, _ := ctx.ParamInt64(EndpointIdParam)
	if rsvalid.IsZero(endpointId) {
		return amerr.GetErrorsFromCode(amerr.ErrEndpointNotFound).GetErrFromLanguage(lang)
	}

	endpoint := &models.Endpoint{Id: endpointId}
	if err := handler.endpointService.DeleteEndpointById(endpoint); err != nil {
		return err.GetErrFromLanguage(lang)
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (handler *EndpointHandlerImpl) GetEndpointList(c echo.Context) error {
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

	list, aerr := handler.endpointService.GetEndpointList(models.EndpointListRequest{
		Page:         int(pageInt64),
		NumItem:      int(numItemInt64),
		WebServiceId: 0,
	})
	if aerr != nil {
		return aerr.GetErrFromLanguage(lang)
	}

	return ctx.JSON(http.StatusOK, list)
}

func NewEndpointHandler(webServiceService services.WebServiceService, endpointService services.EndpointService) (EndpointHandler, error) {
	if rsvalid.IsZero(
		webServiceService,
		endpointService,
	) {
		return nil, errors.Wrap(rserrors.ErrInvalidParameter, "EndpointHandler")
	}
	return &EndpointHandlerImpl{
		webServiceService: webServiceService,
		endpointService:   endpointService,
	}, nil
}
