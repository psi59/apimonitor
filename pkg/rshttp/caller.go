package rshttp

import (
	"io/ioutil"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func Get(request Request) (Response, error) {
	req, err := http.NewRequest(http.MethodGet, request.GetUrl(), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req.Header = request.Header.GetHttpHeader()

	return executeRequest(req)
}

func Post(request Request) (Response, error) {
	panic("not implement")
}

func executeRequest(req *http.Request) (Response, error) {
	start := time.Now()
	res, err := http.DefaultClient.Do(req)
	end := time.Since(start)

	result, err := getResponse(res, end.Milliseconds())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func getResponse(res *http.Response, responseTime int64) (Response, error) {
	rawBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	body := make(map[string]interface{})
	if err := jsoniter.Unmarshal(rawBody, &body); err != nil {
		return nil, errors.WithStack(err)
	}

	return &HttpResponse{
		StatusCode:   res.StatusCode,
		ResponseTime: responseTime,
		Body:         body,
	}, nil
}
