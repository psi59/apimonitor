package models

import (
	"database/sql/driver"
	"encoding/json"
	"net/url"
	"time"

	"github.com/pkg/errors"

	"github.com/realsangil/apimonitor/pkg/rsdb"
	"github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/rshttp"
	"github.com/realsangil/apimonitor/pkg/rsjson"
	"github.com/realsangil/apimonitor/pkg/rslog"
	"github.com/realsangil/apimonitor/pkg/rsmodels"
	"github.com/realsangil/apimonitor/pkg/rsvalid"
)

type Test struct {
	rsmodels.DefaultValidateChecker
	Id           int64               `json:"id"`
	WebServiceId int64               `json:"-"`
	Path         rshttp.EndpointPath `json:"path"`
	HttpMethod   rshttp.Method       `json:"http_method"`
	ContentType  rshttp.ContentType  `json:"content_type"`
	Desc         string              `json:"desc" gorm:"Type:TEXT"`
	RequestData  rsjson.MapJson      `json:"request_data" gorm:"Type:JSON"`
	Header       rshttp.Header       `json:"header" gorm:"Type:JSON"`
	QueryParam   rshttp.Query        `json:"query_param" gorm:"Type:JSON"`
	Timeout      rshttp.Timeout      `json:"timeout"`
	Assertion    AssertionV1         `json:"assertion" gorm:"Type:JSON"`
	Created      time.Time           `json:"created"`
	LastModified time.Time           `json:"last_modified"`
}

func (test *Test) UpdateFromRequest(request TestRequest) error {
	test.Path = request.Path
	test.HttpMethod = request.HttpMethod
	test.ContentType = request.ContentType
	test.Desc = request.Desc
	test.RequestData = request.RequestData
	test.Header = request.Header
	test.QueryParam = request.QueryParam
	test.LastModified = time.Now()
	return test.Validate()
}

func (test *Test) Validate() error {
	if rsvalid.IsZero(
		test.WebServiceId,
		test.Path,
		test.HttpMethod,
		test.ContentType,
		test.Created,
		test.LastModified,
	) {
		return errors.Wrap(rserrors.ErrInvalidParameter, "Endpoint")
	}
	if err := test.Path.Validate(); err != nil {
		return errors.WithStack(err)
	}
	if err := test.HttpMethod.Validate(); err != nil {
		return errors.WithStack(err)
	}
	if err := test.ContentType.Validate(); err != nil {
		return err
	}
	test.SetValidated()
	return nil
}

func (test Test) Execute(webService *WebService) (rshttp.Response, error) {
	request, err := test.ToHttpRequest(webService)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err := rshttp.Do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	rslog.Debugf("res='%+v'", res)

	return res, nil
}

func (test Test) ToHttpRequest(webService *WebService) (*rshttp.Request, error) {
	if rsvalid.IsZero(webService) {
		return nil, errors.WithStack(rserrors.ErrInvalidParameter)
	}

	rawUrl := url.URL{
		Scheme: webService.HttpSchema,
		Host:   webService.Host,
		Path:   test.Path.String(),
	}
	rslog.Debugf("rawUrl='%s'", rawUrl.String())

	request := rshttp.Request{
		RawUrl:  rawUrl.String(),
		Header:  test.Header,
		Query:   test.QueryParam,
		Body:    test.RequestData,
		Timeout: test.Timeout,
	}

	return &request, nil
}

func NewTest(webService *WebService, request TestRequest) (*Test, error) {
	if rsvalid.IsZero(webService, request) {
		return nil, errors.Wrap(rserrors.ErrInvalidParameter, "Endpoint")
	}
	test := &Test{
		WebServiceId: webService.Id,
		Created:      time.Now(),
	}
	if err := test.UpdateFromRequest(request); err != nil {
		return nil, errors.WithStack(err)
	}
	return test, nil
}

type TestRequest struct {
	Path        rshttp.EndpointPath `json:"path"`
	HttpMethod  rshttp.Method       `json:"http_method"`
	ContentType rshttp.ContentType  `json:"content_type"`
	Desc        string              `json:"desc"`
	RequestData rsjson.MapJson      `json:"request_data" gorm:"Type:JSON"`
	Header      rshttp.Header       `json:"header" gorm:"Type:JSON"`
	QueryParam  rshttp.Query        `json:"query_param" gorm:"Type:JSON"`
}

func (request TestRequest) Validate() error {
	if rsvalid.IsZero(
		request.Path,
		request.HttpMethod,
		request.ContentType,
	) {
		return errors.Wrap(rserrors.ErrInvalidParameter, "EndpointRequest")
	}
	if err := request.Path.Validate(); err != nil {
		return errors.WithStack(err)
	}
	if err := request.HttpMethod.Validate(); err != nil {
		return errors.WithStack(err)
	}
	if err := request.ContentType.Validate(); err != nil {
		return err
	}
	return nil
}

type TestListItem struct {
	Id           int64               `json:"id"`
	WebServiceId int64               `json:"-"`
	WebService   *WebService         `json:"web_service" gorm:"foreignkey:WebServiceId"`
	Path         rshttp.EndpointPath `json:"path"`
	HttpMethod   rshttp.Method       `json:"http_method"`
	Desc         string              `json:"desc"`
	Created      time.Time           `json:"created"`
	LastModified time.Time           `json:"last_modified"`
}

func (testListItem TestListItem) MarshalJSON() ([]byte, error) {
	endpointUrl := &url.URL{
		Scheme: testListItem.WebService.HttpSchema,
		Host:   testListItem.WebService.Host,
		Path:   testListItem.Path.String(),
	}
	return json.Marshal(struct {
		Id           int64               `json:"id"`
		Path         rshttp.EndpointPath `json:"path"`
		Url          string              `json:"url"`
		HttpMethod   rshttp.Method       `json:"http_method"`
		Desc         string              `json:"desc"`
		Created      time.Time           `json:"created"`
		LastModified time.Time           `json:"last_modified"`
	}{
		Id:           testListItem.Id,
		Path:         testListItem.Path,
		Url:          endpointUrl.String(),
		HttpMethod:   testListItem.HttpMethod,
		Desc:         testListItem.Desc,
		Created:      testListItem.Created,
		LastModified: testListItem.LastModified,
	})
}

func (testListItem TestListItem) TableName() string {
	return "tests"
}

type TestListRequest struct {
	Page         int
	NumItem      int
	WebServiceId int64
}

type AssertionV1 struct {
	StatusCode int
}

func (assertion *AssertionV1) Scan(src interface{}) error {
	return rsdb.ScanJson(assertion, src)
}

func (assertion AssertionV1) Value() (driver.Value, error) {
	return rsdb.JsonValue(assertion)
}

func (assertion AssertionV1) Assert(res rshttp.Response) bool {
	return !rsvalid.IsZero(res) && assertion.StatusCode == res.GetStatusCode()
}
