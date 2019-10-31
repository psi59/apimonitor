package models

import (
	"regexp"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/rsmodels"
	"github.com/realsangil/apimonitor/pkg/rsvalid"
)

var regexExtractHost = regexp.MustCompile(`^(?:(?:(https?)?(?:\:?\/\/))|(?:\/\/))?(((?:\w{1,100}\.)?\w{2,300}\.\w{2,100})(\.\w{2,100})*)`)

type WebService struct {
	rsmodels.DefaultValidateChecker
	Id           int64              `json:"id" gorm:"private_key"`
	Host         string             `json:"host" gorm:"unique"`
	HttpSchema   string             `json:"http_schema" gorm:"Size:20;Default:'http'"`
	Desc         string             `json:"desc" gorm:"Type:TEXT"`
	Favicon      string             `json:"favicon" gorm:"Type:Text"`
	Schedule     WebServiceSchedule `json:"schedule" gorm:"Size:20"`
	Tests        []WebServiceTest   `json:"-" gorm:"foreignkey:WebServiceId"`
	Created      time.Time          `json:"created"`
	LastModified time.Time          `json:"last_modified"`
}

func (webService *WebService) Validate() error {
	if rsvalid.IsZero(
		webService.LastModified, webService.Created, webService.Host,
	) {
		return errors.Wrap(rserrors.ErrInvalidParameter, "webService")
	}
	webService.SetValidated()
	return nil
}

func (webService *WebService) UpdateFromRequest(request WebServiceRequest) error {
	host, err := hostRegexpFindStringSubmatch(request.Host)
	if err != nil {
		return errors.WithStack(err)
	}

	webService.Host = host[2]
	webService.HttpSchema = host[1]
	webService.Desc = request.Desc
	webService.Favicon = request.Favicon
	webService.LastModified = time.Now()

	return nil
}

func (webService *WebService) GetScheduleTicker() time.Ticker {
	panic("not implement")
}

func NewWebService(request WebServiceRequest) (*WebService, error) {
	webService := &WebService{
		Created: time.Now(),
	}

	if err := webService.UpdateFromRequest(request); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := webService.Validate(); err != nil {
		return nil, errors.WithStack(err)
	}

	return webService, nil
}

func hostRegexpFindStringSubmatch(host string) ([]string, error) {
	if !regexExtractHost.MatchString(host) {
		return nil, rserrors.ErrInvalidParameter
	}
	return regexExtractHost.FindStringSubmatch(host), nil
}

type WebServiceRequest struct {
	Host    string `json:"host"`
	Desc    string `json:"desc"`
	Favicon string `json:"favicon"`
}

type WebServiceListRequest struct {
	Page          int
	NumItem       int
	SearchKeyword string
}

const (
	ScheduleOneMinute     WebServiceSchedule = "1m"
	ScheduleFiveMinute    WebServiceSchedule = "5m"
	ScheduleFifteenMinute WebServiceSchedule = "15m"
	ScheduleThirtyMinute  WebServiceSchedule = "30m"
	ScheduleHourly        WebServiceSchedule = "1h"
	ScheduleDaily         WebServiceSchedule = "1d"
)

type WebServiceSchedule string

func (schedule *WebServiceSchedule) Validate() error {
	switch *schedule {
	case ScheduleOneMinute, ScheduleFiveMinute, ScheduleFifteenMinute, ScheduleThirtyMinute, ScheduleHourly, ScheduleDaily:
		return nil
	default:
		return rserrors.Error("invalid schedule")
	}
}

func (schedule *WebServiceSchedule) UnmarshalJSON(data []byte) error {
	var str string
	if err := jsoniter.Unmarshal(data, &str); err != nil {
		return errors.WithStack(err)
	}
	s := WebServiceSchedule(str)
	if str == "" {
		s = ScheduleDaily
	}
	if err := s.Validate(); err != nil {
		return errors.WithStack(err)
	}
	*schedule = s
	return schedule.Validate()
}

func (schedule WebServiceSchedule) GetDuration() time.Duration {
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

func (schedule WebServiceSchedule) GetTicker() *time.Ticker {
	return time.NewTicker(schedule.GetDuration())
}
