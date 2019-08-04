package rsjson

import (
	"database/sql/driver"
	"errors"

	jsoniter "github.com/json-iterator/go"
)

type MapJson map[string]interface{}

func (mapJson *MapJson) Value() (driver.Value, error) {
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
