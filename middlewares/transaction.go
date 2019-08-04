package middlewares

import (
	"github.com/labstack/echo"
	"github.com/pkg/errors"

	"github.com/realsangil/apimonitor/pkg/rsdb"
	"github.com/realsangil/apimonitor/pkg/rslog"
)

type TransactionMiddleware interface {
	Tx(handlerFunc echo.HandlerFunc) echo.HandlerFunc
}

type TransactionMiddlewareImpl struct{}

func (middleware *TransactionMiddlewareImpl) Tx(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tx, err := rsdb.GetConnection().Begin()
		if err != nil {
			return errors.WithStack(err)
		}
		defer func() {
			if err := tx.Rollback(); err != nil {
				rslog.Error(err)
			}
		}()

		ctx, err := ConvertToCustomContext(c)
		if err != nil {
			return errors.WithStack(err)
		}
		if err := ctx.SetTx(tx); err != nil {
			return errors.WithStack(err)
		}

		if err := next(ctx); err != nil {
			return errors.WithStack(err)
		}

		if err := tx.Commit(); err != nil {
			return errors.WithStack(err)
		}

		return nil
	}
}

func NewTransactionMiddleware() TransactionMiddleware {
	return &TransactionMiddlewareImpl{}
}
