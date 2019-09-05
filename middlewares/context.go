package middlewares

import (
	"strconv"

	"github.com/realsangil/apimonitor/pkg/rsdb"
	"github.com/realsangil/apimonitor/pkg/rserrors"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

const ErrInvalidContext = rserrors.Error("invalid context")

type Context interface {
	echo.Context
	SetTx(transaction rsdb.Connection) error
	GetTx() (rsdb.Connection, error)
	Language() string
	QueryParamInt64(name string, defaultInt int64) (int64, error)
	ParamInt64(name string) (int64, error)
}

type context struct {
	echo.Context
}

func (c *context) SetTx(transaction rsdb.Connection) error {
	if transaction == nil {
		return rsdb.ErrInvalidTransaction
	}
	c.Set("tx", transaction)
	return nil
}

func (c *context) GetTx() (rsdb.Connection, error) {
	tx, ok := c.Get("tx").(rsdb.Connection)
	if !ok {
		return nil, rsdb.ErrInvalidTransaction
	}
	return tx, nil
}

func (c *context) JSON(code int, data interface{}) error {
	type Result struct {
		Success bool        `json:"success"`
		Result  interface{} `json:"result,omitempty"`
	}
	return c.Context.JSON(code, Result{
		Success: true,
		Result:  data,
	})
}

func (c *context) Language() string {
	lang, exists := c.Request().Header["accept-language"]
	if !exists {
		lang = []string{"ko"}
	}
	return lang[0]
}

func (c *context) QueryParamInt64(name string, defaultValue int64) (int64, error) {
	rawQueryParam := c.QueryParam(name)
	return c.stringToInt64WithDefault(rawQueryParam, defaultValue)
}

func (c *context) ParamInt64(name string) (int64, error) {
	rawParam := c.Param(name)
	return c.stringToInt64WithDefault(rawParam, 0)
}

func (c *context) stringToInt64WithDefault(value string, defaultValue int64) (int64, error) {
	if value == "" {
		return defaultValue, nil
	}
	return strconv.ParseInt(value, 10, 64)
}

func ConvertToCustomContext(c echo.Context) (Context, error) {
	ctx, ok := c.(Context)
	if !ok {
		return nil, errors.WithStack(ErrInvalidContext)
	}
	return ctx, nil
}

func NewContext(c echo.Context) Context {
	return &context{
		Context: c,
	}
}

func ReplaceContextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := NewContext(c)
		return next(cc)
	}
}
