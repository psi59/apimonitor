package rsjson

import (
	"database/sql/driver"
	"errors"
	"reflect"

	jsoniter "github.com/json-iterator/go"
)

type MapJson map[string]interface{}

func (mapJson MapJson) Value() (driver.Value, error) {
	return jsoniter.Marshal(mapJson)
}

func (mapJson *MapJson) Scan(src interface{}) error {
	if src == nil {
		*mapJson = nil
		return nil
	}
	s, ok := src.([]byte)
	if !ok {
		return errors.New("Invalid Scan Source")
	}
	return jsoniter.Unmarshal(s, mapJson)
}

func ValueToDB(src interface{}) (driver.Value, error) {
	return jsoniter.Marshal(src)
}

func ScanFromDB(src, dst interface{}) error {
	dstType := reflect.TypeOf(dst)
	if dstType.Kind() != reflect.Ptr {
		return errors.New("dst is not ptr")
	}

	if src == nil {
		dst = nil
		return nil
	}

	s, ok := src.([]byte)
	if !ok {
		return errors.New("invalid scan source")
	}
	return jsoniter.Unmarshal(s, dst)
}
