package rshttp

import (
	"net/http"
	"time"

	"github.com/imroc/req"
	"github.com/realsangil/apimonitor/pkg/rsjson"
	"github.com/realsangil/apimonitor/pkg/rslog"
	"github.com/realsangil/apimonitor/pkg/rsvalid"
)

const (
	DefaultTimeout = 5 * time.Second
)

type Timeout int64

type requestFunc func(url string, v ...interface{}) (*req.Resp, error)

func (timeout Timeout) GetDuration() time.Duration {
	if rsvalid.IsZero(timeout) {
		return DefaultTimeout
	}
	return time.Duration(timeout) * time.Second
}

type Request struct {
	Method      string
	RawUrl      string
	Header      req.Header
	Query       req.QueryParam
	ContentType ContentType
	Body        rsjson.MapJson
	Timeout     Timeout
}

func (request *Request) execute() (*Response, error) {
	var requestFn requestFunc
	switch request.Method {
	case http.MethodGet:
		requestFn = req.Get
	case http.MethodPost:
		requestFn = req.Post
	case http.MethodPut:
		requestFn = req.Put
	case http.MethodDelete:
		requestFn = req.Delete
	default:
		return nil, ErrUnsupportedMethod
	}
	return execute(requestFn, request)
}

func execute(fn requestFunc, request *Request) (*Response, error) {
	resp, err := fn(request.RawUrl, request.Header, request.Query, req.BodyJSON(request.Body))
	if err != nil {
		rslog.Error(err)
		return nil, err
	}
	return &Response{
		StatusCode:   resp.Response().StatusCode,
		ResponseTime: resp.Cost().Milliseconds(),
		Body:         resp.String(),
	}, nil
}

type Response struct {
	StatusCode   int
	ResponseTime int64
	Body         string
}
