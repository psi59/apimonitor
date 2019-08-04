package testutils

import (
	"time"

	"bou.ke/monkey"
	"github.com/pkg/errors"
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

func MonkeyAll() {
	monkey.UnpatchAll()
	MonkeyErrorsWrap()
	MonkeyErrorsWithStack()
	MonkeyTimeNow(time.Now())
}
