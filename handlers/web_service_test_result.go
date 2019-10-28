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

type WebServiceTestResultHandler interface {
	GetList(c echo.Context) error
}

type WebServiceTestResultHandlerImpl struct {
	webServiceTestResultService services.WebServiceTestResultService
}

func (handler *WebServiceTestResultHandlerImpl) GetList(c echo.Context) error {
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

	webServiceTestId, err := ctx.ParamInt64(WebServiceTestIdParam)
	if err != nil {
		rslog.Error(err)
		return amerr.GetErrorsFromCode(amerr.ErrWebServiceTestNotFound).GetErrFromLanguage(lang)
	}

	request := models.WebServiceTestResultListRequest{
		Page:      int(page),
		NumItem:   int(numItem),
		IsSuccess: models.IsSuccess(isSuccess),
		// StartTestedAt: time.Time{},
		// EndTestedAt:   time.Time{},
	}
	webServiceTest := &models.WebServiceTest{Id: webServiceTestId}

	list, aerr := handler.webServiceTestResultService.GetResultListByTestId(webServiceTest, request)
	if aerr != nil {
		return aerr.GetErrFromLanguage(lang)
	}

	return ctx.JSON(http.StatusOK, list)
}

func NewWebServiceTestResultHandler(webServiceTestResultService services.WebServiceTestResultService) (WebServiceTestResultHandler, error) {
	if rsvalid.IsZero(webServiceTestResultService) {
		return nil, rserrors.ErrInvalidParameter
	}
	return &WebServiceTestResultHandlerImpl{webServiceTestResultService: webServiceTestResultService}, nil
}
