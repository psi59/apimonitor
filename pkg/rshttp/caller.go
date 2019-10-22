package rshttp

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/rsvalid"
)

func Do(req *Request) (Response, error) {
	if rsvalid.IsZero(req) {
		return nil, errors.WithStack(rserrors.ErrInvalidParameter)
	}
	rawBody, _ := jsoniter.Marshal(req.Body)
	var requestBody io.Reader
	if !rsvalid.IsZero(rawBody) {
		requestBody = bytes.NewReader(rawBody)
	}
	ctx, cancel := context.WithTimeout(context.Background(), req.Timeout.GetDuration())

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, req.GetUrl(), requestBody)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	httpRequest.Header = req.Header.GetHttpHeader()

	errChan := make(chan error, 1)

	httpResponse := new(http.Response)
	getResponse := func(httpResponse *http.Response, endDuration time.Duration, err error) (Response, error) {
		if err != nil {
			return nil, err
		}

		// TODO: response 바디 가져오기 추가
		// result, err := responseContentType.GetBodyFromResponse(res)
		// if err != nil {
		// 	return nil, errors.WithStack(err)
		// }

		return &HttpResponse{
			StatusCode:   httpResponse.StatusCode,
			ResponseTime: endDuration.Milliseconds(),
			Body:         nil,
		}, nil
	}
	start := time.Now()
	go func() {
		defer close(errChan)
		res, err := http.DefaultClient.Do(httpRequest)
		if err != nil {
			errChan <- err
		}
		httpResponse = res
		errChan <- nil
	}()
	end := time.Since(start)

	select {
	case <-ctx.Done():
		cancel()
		return nil, rserrors.Error("timeout")
	case err := <-errChan:
		return getResponse(httpResponse, end, err)
	}
}
