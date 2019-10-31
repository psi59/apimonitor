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

type WebServiceTest struct {
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

func (webServiceTest *WebServiceTest) UpdateFromRequest(request WebServiceTestRequest) error {
	webServiceTest.Path = request.Path
	webServiceTest.HttpMethod = request.HttpMethod
	webServiceTest.ContentType = request.ContentType
	webServiceTest.Desc = request.Desc
	webServiceTest.RequestData = request.RequestData
	webServiceTest.Header = request.Header
	webServiceTest.QueryParam = request.QueryParam
	webServiceTest.LastModified = time.Now()
	return webServiceTest.Validate()
}

func (webServiceTest *WebServiceTest) Validate() error {
	if rsvalid.IsZero(
		webServiceTest.WebServiceId,
		webServiceTest.Path,
		webServiceTest.HttpMethod,
		webServiceTest.ContentType,
		webServiceTest.Created,
		webServiceTest.LastModified,
	) {
		return errors.Wrap(rserrors.ErrInvalidParameter, "Endpoint")
	}
	if err := webServiceTest.Path.Validate(); err != nil {
		return errors.WithStack(err)
	}
	if err := webServiceTest.HttpMethod.Validate(); err != nil {
		return errors.WithStack(err)
	}
	if err := webServiceTest.ContentType.Validate(); err != nil {
		return err
	}
	webServiceTest.SetValidated()
	return nil
}

func (webServiceTest WebServiceTest) Execute(webService *WebService) (rshttp.Response, error) {
	request, err := webServiceTest.ToHttpRequest(webService)
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

func (webServiceTest WebServiceTest) ToHttpRequest(webService *WebService) (*rshttp.Request, error) {
	if rsvalid.IsZero(webService) {
		return nil, errors.WithStack(rserrors.ErrInvalidParameter)
	}

	rawUrl := url.URL{
		Scheme: webService.HttpSchema,
		Host:   webService.Host,
		Path:   webServiceTest.Path.String(),
	}
	rslog.Debugf("rawUrl='%s'", rawUrl.String())

	request := rshttp.Request{
		RawUrl:  rawUrl.String(),
		Header:  webServiceTest.Header,
		Query:   webServiceTest.QueryParam,
		Body:    webServiceTest.RequestData,
		Timeout: webServiceTest.Timeout,
	}

	return &request, nil
}

func NewWebServiceTest(webService *WebService, request WebServiceTestRequest) (*WebServiceTest, error) {
	if rsvalid.IsZero(webService, request) {
		return nil, errors.Wrap(rserrors.ErrInvalidParameter, "Endpoint")
	}
	webServiceTest := &WebServiceTest{
		WebServiceId: webService.Id,
		Created:      time.Now(),
	}
	if err := webServiceTest.UpdateFromRequest(request); err != nil {
		return nil, errors.WithStack(err)
	}
	return webServiceTest, nil
}

type WebServiceTestRequest struct {
	Path        rshttp.EndpointPath `json:"path"`
	HttpMethod  rshttp.Method       `json:"http_method"`
	ContentType rshttp.ContentType  `json:"content_type"`
	Desc        string              `json:"desc"`
	RequestData rsjson.MapJson      `json:"request_data" gorm:"Type:JSON"`
	Header      rshttp.Header       `json:"header" gorm:"Type:JSON"`
	QueryParam  rshttp.Query        `json:"query_param" gorm:"Type:JSON"`
}

func (request WebServiceTestRequest) Validate() error {
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

type WebServiceTestListItem struct {
	Id           int64               `json:"id"`
	WebServiceId int64               `json:"-"`
	WebService   *WebService         `json:"web_service" gorm:"foreignkey:WebServiceId"`
	Path         rshttp.EndpointPath `json:"path"`
	HttpMethod   rshttp.Method       `json:"http_method"`
	Desc         string              `json:"desc"`
	Created      time.Time           `json:"created"`
	LastModified time.Time           `json:"last_modified"`
}

func (webServiceTestListItem WebServiceTestListItem) MarshalJSON() ([]byte, error) {
	endpointUrl := &url.URL{
		Scheme: webServiceTestListItem.WebService.HttpSchema,
		Host:   webServiceTestListItem.WebService.Host,
		Path:   webServiceTestListItem.Path.String(),
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
		Id:           webServiceTestListItem.Id,
		Path:         webServiceTestListItem.Path,
		Url:          endpointUrl.String(),
		HttpMethod:   webServiceTestListItem.HttpMethod,
		Desc:         webServiceTestListItem.Desc,
		Created:      webServiceTestListItem.Created,
		LastModified: webServiceTestListItem.LastModified,
	})
}

func (webServiceTestListItem WebServiceTestListItem) TableName() string {
	return "web_service_tests"
}

type WebServiceTestListRequest struct {
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
