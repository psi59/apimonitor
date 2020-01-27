package models

import (
	"database/sql/driver"

	"github.com/imroc/req"
	jsoniter "github.com/json-iterator/go"
	"github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/rslog"
)

type Alerter interface {
	Alert(msg string) error
}

type WebHookAlerts []*WebHookAlerter

func (alerts *WebHookAlerts) Scan(src interface{}) error {
	data, ok := src.([]byte)
	if !ok {
		return rserrors.Error("invalid data")
	}
	return jsoniter.Unmarshal(data, alerts)
}

func (alerts WebHookAlerts) Value() (driver.Value, error) {
	return jsoniter.Marshal(alerts)
}

type WebHookAlerter struct {
	URL    string `json:"url"`
	Enable bool   `json:"disable"`
}

func (alerter WebHookAlerter) Alert(msg string) error {
	resp, err := req.Post(alerter.URL, req.BodyJSON(map[string]string{
		"text": msg,
	}))
	rslog.Debug(resp)
	return err
}
