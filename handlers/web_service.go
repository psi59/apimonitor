package handlers

import (
	"net/http"

	"github.com/labstack/echo"
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
	WebServiceIdParam = "web_service_id"
)

type WebServiceHandler interface {
	CreateWebService(c echo.Context) error
	GetWebServiceById(c echo.Context) error
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
	tx, err := ctx.GetTx()
	if err != nil {
		rslog.Error(err)
		return amerr.GetErrInternalServer().GetErrFromLanguage(lang)
	}

	var request models.WebServiceRequest
	if err := ctx.Bind(&request); err != nil {
		rslog.Error(err)
		return amerr.GetErrorsFromCode(amerr.ErrBadRequest).GetErrFromLanguage(lang)
	}
	webService, aerr := handler.webServiceService.CreateWebService(tx, request)
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
	tx, err := ctx.GetTx()
	if err != nil {
		rslog.Error(err)
		return amerr.GetErrInternalServer().GetErrFromLanguage(lang)
	}

	webServiceId, err := ctx.ParamInt64(WebServiceIdParam)
	if err != nil {
		rslog.Error(err)
		return amerr.GetErrorsFromCode(amerr.ErrWebServiceNotFound).GetErrFromLanguage(lang)
	}

	webService := &models.WebService{Id: webServiceId}
	if err := handler.webServiceService.GetWebServiceById(tx, webService); err != nil {
		return err.GetErrFromLanguage(lang)
	}

	return ctx.JSON(http.StatusOK, webService)
}

func NewWebServiceHandler(webServiceService services.WebServiceService) (WebServiceHandler, error) {
	if rsvalid.IsZero(webServiceService) {
		return nil, rserrors.ErrInvalidParameter
	}
	return &WebServiceHandlerImpl{
		webServiceService: webServiceService,
	}, nil
}
