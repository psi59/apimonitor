package rshttp

import (
	"net/http"

	"github.com/imroc/req"
	"github.com/pkg/errors"

	"github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/rslog"
	"github.com/realsangil/apimonitor/pkg/rsvalid"
)

func Do(request *Request) (*Response, error) {
	if rsvalid.IsZero(request) {
		return nil, errors.WithStack(rserrors.ErrInvalidParameter)
	}

	switch request.Method {
	case http.MethodGet:
		resp, err := req.Get(request.GetUrl())
		if err != nil {
			rslog.Error(err)
			return nil, err
		}

		return &Response{
			StatusCode:   resp.Response().StatusCode,
			ResponseTime: resp.Cost().Milliseconds(),
			Body:         resp.String(),
		}, nil
	default:
		return nil, rserrors.ErrUnexpected
	}
	// rawBody, _ := jsoniter.Marshal(request.Body)
	// var requestBody io.Reader
	// if !rsvalid.IsZero(rawBody) {
	// 	requestBody = bytes.NewReader(rawBody)
	// }
	// ctx, cancel := context.WithTimeout(context.Background(), request.Timeout.GetDuration())
	//
	// httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, request.GetUrl(), requestBody)
	// if err != nil {
	// 	return nil, errors.WithStack(err)
	// }
	// httpRequest.Header = request.Header.GetHttpHeader()
	//
	// errChan := make(chan error, 1)
	//
	// httpResponse := new(http.Response)
	// getResponse := func(httpResponse *http.Response, endDuration time.Duration, err error) (Response, error) {
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	// 	// TODO: response 바디 가져오기 추가
	// 	// result, err := responseContentType.GetBodyFromResponse(res)
	// 	// if err != nil {
	// 	// 	return nil, errors.WithStack(err)
	// 	// }
	//
	// 	return &Response{
	// 		StatusCode:   httpResponse.StatusCode,
	// 		ResponseTime: endDuration.Milliseconds(),
	// 		Body:         nil,
	// 	}, nil
	// }
	// start := time.Now()
	// go func() {
	// 	defer close(errChan)
	// 	res, err := http.DefaultClient.Do(httpRequest)
	// 	if err != nil {
	// 		errChan <- err
	// 	}
	// 	httpResponse = res
	// 	errChan <- nil
	// }()
	// end := time.Since(start)
	//
	// select {
	// case <-ctx.Done():
	// 	cancel()
	// 	return nil, rserrors.Error("timeout")
	// case err := <-errChan:
	// 	return getResponse(httpResponse, end, err)
	// }
}
