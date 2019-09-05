package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/realsangil/apimonitor/pkg/amerr"
	"github.com/realsangil/apimonitor/pkg/rslog"
)

func ErrorHandleMiddleware(err error, c echo.Context) {
	statusCode := http.StatusInternalServerError
	all := err.Error()
	code := 0
	switch e := errors.Cause(err).(type) {
	case *echo.HTTPError:
		statusCode = e.Code
		all = e.Message.(string)

	case *amerr.Error:
		all = e.Message
		statusCode = e.StatusCode
		code = e.ErrorCode
	}

	if statusCode == http.StatusInternalServerError {
		all = "Internal Server Error"
		rslog.Errorf("%+v\n", err)
	}

	result := map[string]interface{}{
		"sucess": false,
		"all":    all,
		"code":   code,
	}

	c.JSON(statusCode, result)
}
