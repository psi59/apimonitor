package models

import (
	"database/sql/driver"
	"encoding/json"
	"net/url"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/realsangil/apimonitor/pkg/rsdb"
	"github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/rshttp"
	"github.com/realsangil/apimonitor/pkg/rsjson"
	"github.com/realsangil/apimonitor/pkg/rslog"
	"github.com/realsangil/apimonitor/pkg/rsmodels"
	"github.com/realsangil/apimonitor/pkg/rsstr"
	"github.com/realsangil/apimonitor/pkg/rsvalid"
)

type Test struct {
	rsmodels.DefaultValidateChecker
	Id           string              `json:"id"`
	Name         string              `json:"name"`
	WebServiceId string              `json:"-"`
	WebService   WebService          `json:"webService,omitempty" gorm:"PRELOAD:true"`
	Path         rshttp.EndpointPath `json:"path"`
	Method       rshttp.Method       `json:"method"`
	ContentType  rshttp.ContentType  `json:"contentType"`
	Description  string              `json:"description" gorm:"Type:TEXT"`
	Parameters   Parameters          `json:"parameters" gorm:"Type:JSON"`
	Schedule     TestSchedule        `json:"schedule" gorm:"Column:schedule;Type:VARCHAR(5)"`
	Timeout      rshttp.Timeout      `json:"timeout"`
	Assertion    AssertionV1         `json:"assertion" gorm:"Type:JSON"`
	CreatedAt    time.Time           `json:"createdAt"`
	ModifiedAt   time.Time           `json:"modifiedAt"`
}

func (test *Test) UpdateFromRequest(request TestRequest) error {
	test.Name = request.Name
	test.Path = request.Path
	test.Method = request.Method
	test.ContentType = request.ContentType
	test.Description = request.Description
	test.Schedule = request.Schedule
	test.Timeout = rshttp.Timeout(request.Timeout)
	test.ModifiedAt = time.Now()
	return test.Validate()
}

func (test *Test) Validate() error {
	if rsvalid.IsZero(
		test.Id,
		test.Name,
		test.WebServiceId,
		test.Path,
		test.Method,
		test.ContentType,
		test.Schedule,
		test.CreatedAt,
		test.ModifiedAt,
	) {
		return errors.Wrap(rserrors.ErrInvalidParameter, "Endpoint")
	}
	if err := test.Path.Validate(); err != nil {
		return errors.WithStack(err)
	}
	if err := test.Method.Validate(); err != nil {
		return errors.WithStack(err)
	}
	if err := test.ContentType.Validate(); err != nil {
		return errors.WithStack(err)
	}
	if err := test.Schedule.Validate(); err != nil {
		return errors.WithStack(err)
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
		Scheme: webService.Schema,
		Host:   webService.Host,
		Path:   test.Path.String(),
	}
	rslog.Debugf("rawUrl='%s'", rawUrl.String())

	request := rshttp.Request{
		RawUrl:  rawUrl.String(),
		Timeout: test.Timeout,
	}

	return &request, nil
}

func NewTest(webService *WebService, request TestRequest) (*Test, error) {
	if rsvalid.IsZero(webService, request) {
		return nil, errors.Wrap(rserrors.ErrInvalidParameter, "Endpoint")
	}
	test := &Test{
		Id:           rsstr.NewUUID(),
		WebServiceId: webService.Id,
		CreatedAt:    time.Now(),
	}
	if err := test.UpdateFromRequest(request); err != nil {
		return nil, errors.WithStack(err)
	}
	return test, nil
}

type TestRequest struct {
	Id          string              `json:"-"`
	Name        string              `json:"name"`
	Path        rshttp.EndpointPath `json:"path"`
	Method      rshttp.Method       `json:"method"`
	ContentType rshttp.ContentType  `json:"contentType"`
	Description string              `json:"description"`
	Parameters  Parameters          `json:"parameters"`
	Schedule    TestSchedule        `json:"schedule"`
	Assertion   AssertionV1         `json:"assertion"`
	Timeout     int                 `json:"timeout"`
}

func (request TestRequest) Validate() error {
	if rsvalid.IsZero(
		request.Path,
		request.Method,
		request.ContentType,
	) {
		return errors.Wrap(rserrors.ErrInvalidParameter, "EndpointRequest")
	}
	if err := request.Path.Validate(); err != nil {
		return errors.WithStack(err)
	}
	if err := request.Method.Validate(); err != nil {
		return errors.WithStack(err)
	}
	if err := request.ContentType.Validate(); err != nil {
		return err
	}
	return nil
}

type TestListItem struct {
	Id           string              `json:"id"`
	WebServiceId string              `json:"webServiceId"  gorm:"foreignkey:WebServiceId"`
	WebService   *WebService         `json:"webService"`
	Path         rshttp.EndpointPath `json:"path"`
	Method       rshttp.Method       `json:"method"`
	Description  string              `json:"description"`
	Created      time.Time           `json:"created"`
	LastModified time.Time           `json:"lastModified"`
}

func (testListItem TestListItem) MarshalJSON() ([]byte, error) {
	endpointUrl := &url.URL{
		Scheme: testListItem.WebService.Schema,
		Host:   testListItem.WebService.Host,
		Path:   testListItem.Path.String(),
	}
	return json.Marshal(struct {
		Id           string              `json:"id"`
		Path         rshttp.EndpointPath `json:"path"`
		Url          string              `json:"url"`
		Method       rshttp.Method       `json:"method"`
		Desc         string              `json:"desc"`
		Created      time.Time           `json:"created"`
		LastModified time.Time           `json:"lastModified"`
	}{
		Id:           testListItem.Id,
		Path:         testListItem.Path,
		Url:          endpointUrl.String(),
		Method:       testListItem.Method,
		Desc:         testListItem.Description,
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
	WebServiceId string
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

const (
	ScheduleOneMinute     TestSchedule = "1m"
	ScheduleFiveMinute    TestSchedule = "5m"
	ScheduleFifteenMinute TestSchedule = "15m"
	ScheduleThirtyMinute  TestSchedule = "30m"
	ScheduleHourly        TestSchedule = "1h"
	ScheduleDaily         TestSchedule = "1d"
)

type TestSchedule string

func (schedule *TestSchedule) Validate() error {
	switch *schedule {
	case ScheduleOneMinute, ScheduleFiveMinute, ScheduleFifteenMinute, ScheduleThirtyMinute, ScheduleHourly, ScheduleDaily:
		return nil
	default:
		return rserrors.Error("invalid schedule")
	}
}

func (schedule *TestSchedule) UnmarshalJSON(data []byte) error {
	var str string
	if err := jsoniter.Unmarshal(data, &str); err != nil {
		return errors.WithStack(err)
	}
	s := TestSchedule(str)
	if str == "" {
		s = ScheduleDaily
	}
	if err := s.Validate(); err != nil {
		return errors.WithStack(err)
	}
	*schedule = s
	return schedule.Validate()
}

func (schedule TestSchedule) GetDuration() time.Duration {
	// return 3 * time.Second
	switch schedule {
	case ScheduleOneMinute:
		return 1 * time.Minute
	case ScheduleFiveMinute:
		return 5 * time.Minute
	case ScheduleFifteenMinute:
		return 15 * time.Minute
	case ScheduleThirtyMinute:
		return 30 * time.Minute
	case ScheduleHourly:
		return time.Hour
	case ScheduleDaily:
		return 24 * time.Hour
	default:
		return 0
	}
}

func (schedule TestSchedule) GetTicker() *time.Ticker {
	return time.NewTicker(schedule.GetDuration())
}

type Parameters struct {
	Auth   map[string]interface{} `json:"auth"`
	Header map[string][]string    `json:"header"`
	Query  map[string][]string    `json:"query"`
	Body   map[string]interface{} `json:"body"`
}

func (parameters *Parameters) Scan(src interface{}) error {
	return rsjson.ScanFromDB(src, parameters)
}

func (parameters Parameters) Value() (driver.Value, error) {
	return jsoniter.Marshal(parameters)
}
