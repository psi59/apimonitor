package rshttp

import (
	"database/sql/driver"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/realsangil/apimonitor/pkg/rsjson"
	"github.com/realsangil/apimonitor/pkg/rslog"
	"github.com/realsangil/apimonitor/pkg/rsvalid"
)

const (
	DefaultTimeout = 5 * time.Second
)

type Header map[string]interface{}

func (header *Header) Scan(src interface{}) error {
	return rsjson.ScanFromDB(src, header)
}

func (header Header) Value() (driver.Value, error) {
	return rsjson.ValueToDB(header)
}

func (header Header) GetHttpHeader() http.Header {
	if rsvalid.IsZero(header) {
		return nil
	}
	var httpHeader http.Header
	for k, v := range header {
		httpHeader.Add(k, convertInterfaceToString(v))
	}

	return httpHeader
}

type Query map[string]interface{}

func (query *Query) Scan(src interface{}) error {
	return rsjson.ScanFromDB(src, query)
}

func (query Query) Value() (driver.Value, error) {
	return rsjson.ValueToDB(query)
}

func (query Query) GetHttpQuery() url.Values {
	if rsvalid.IsZero(query) {
		return nil
	}
	var httpQuery url.Values
	for k, v := range query {
		httpQuery.Add(k, convertInterfaceToString(v))
	}
	return httpQuery
}

func (query Query) GetHttpQueryString() string {
	q := query.GetHttpQuery()
	if q == nil {
		return ""
	}
	return q.Encode()
}

type Timeout int64

func (timeout Timeout) GetDuration() time.Duration {
	if rsvalid.IsZero(timeout) {
		return DefaultTimeout
	}
	return time.Duration(timeout) * time.Second
}

type Request struct {
	RawUrl      string
	Header      Header
	Query       Query
	ContentType ContentType
	Body        rsjson.MapJson
	Timeout     Timeout
}

func (request Request) GetUrl() string {
	parsedUrl, err := url.Parse(request.RawUrl)
	if err != nil {
		rslog.Error(err)
		return ""
	}
	query := request.Query.GetHttpQuery()
	for k, v := range query {
		parsedUrl.Query().Add(k, v[0])
	}
	parsedUrl.RawQuery = parsedUrl.Query().Encode()
	rslog.Infof("request_url='%s'", parsedUrl.String())
	return parsedUrl.String()
}

type Response interface {
	GetStatusCode() int
	GetResponseTime() int64
	GetBody() interface{}
}

type HttpResponse struct {
	StatusCode   int
	ResponseTime int64
	Body         interface{}
}

func (d HttpResponse) GetStatusCode() int {
	return d.StatusCode
}

func (d HttpResponse) GetResponseTime() int64 {
	return d.ResponseTime
}

func (d HttpResponse) GetBody() interface{} {
	return d.Body
}

func convertInterfaceToString(i interface{}) string {
	return fmt.Sprintf("%v", i)
}
