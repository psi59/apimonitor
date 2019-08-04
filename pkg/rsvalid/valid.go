package rsvalid

import (
	"reflect"

	"github.com/realsangil/apimonitor/pkg/rslog"
)

func IsZero(i ...interface{}) bool {
	isZero := false
	for _, j := range i {
		typ := reflect.TypeOf(j)
		if typ == nil {
			return true
		}
		if j == reflect.Zero(typ).Interface() {
			rslog.Debugf("%s is zero", typ)
			isZero = true
		}
	}
	return isZero
}
