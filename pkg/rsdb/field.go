package rsdb

import (
	"database/sql/driver"
	"reflect"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/realsangil/apimonitor/pkg/rserrors"
)

func ScanJson(dst, src interface{}) error {
	if reflect.TypeOf(dst).Kind() != reflect.Ptr {
		return rserrors.Error("dst is not pointer")
	}

	if src == nil {
		dst = nil
		return nil
	}
	s, ok := src.([]byte)
	if !ok {
		return errors.New("Invalid Scan Source")
	}
	return jsoniter.Unmarshal(s, dst)
}

func JsonValue(dst interface{}) (driver.Value, error) {
	return jsoniter.Marshal(dst)
}
