package testutils

import (
	"time"

	"bou.ke/monkey"
	"github.com/pkg/errors"

	"github.com/realsangil/apimonitor/pkg/rsdb"
	"github.com/realsangil/apimonitor/pkg/rsdb/mocks"
)

func MonkeyTimeNow(now time.Time) {
	monkey.Patch(time.Now, func() time.Time {
		return now
	})
}

func MonkeyErrorsWrap() {
	monkey.Patch(errors.Wrap, func(err error, s string) error {
		return err
	})
}

func MonkeyErrorsWithStack() {
	monkey.Patch(errors.WithStack, func(err error) error {
		return err
	})
}

func MonkeyGetConnection(conn *mocks.Connection) {
	monkey.Patch(rsdb.GetConnection, func() rsdb.Connection {
		return conn
	})
}

func MonkeyAll() {
	monkey.UnpatchAll()
	MonkeyErrorsWrap()
	MonkeyErrorsWithStack()
	MonkeyTimeNow(time.Now())
}
