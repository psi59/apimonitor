package services

import (
	"github.com/pkg/errors"

	"github.com/realsangil/apimonitor/models"
	"github.com/realsangil/apimonitor/pkg/rsdb"
	"github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/rslog"
	"github.com/realsangil/apimonitor/pkg/rsvalid"
	"github.com/realsangil/apimonitor/repositories"
)

type ScheduleRunner interface {
	Run() error
}

type ScheduleCloser interface {
	Close() error
}

type ScheduleRefresher interface {
	Refresh() error
}

type WebServiceScheduler interface {
	ScheduleRunner
	ScheduleCloser
}

type WebServiceScheduleManager interface {
	ScheduleRunner
	ScheduleRefresher
	ScheduleCloser
}

type webServiceScheduleManager struct {
	webServiceSchedulers           map[interface{}]WebServiceScheduler
	webServiceRepository           repositories.WebServiceRepository
	webServiceTestResultRepository repositories.WebServiceTestResultRepository
}

func (manager *webServiceScheduleManager) Run() error {
	rslog.Debug("Running WebServiceManager...")
	for _, s := range manager.webServiceSchedulers {
		if err := s.Run(); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (manager *webServiceScheduleManager) Refresh() error {
	if err := manager.Close(); err != nil {
		return errors.WithStack(err)
	}
	if err := manager.refreshWebServices(); err != nil {
		return errors.WithStack(err)
	}
	if err := manager.Run(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (manager *webServiceScheduleManager) refreshWebServices() error {
	rslog.Debug("Refreshing WebServiceScheduler...")
	webServices, err := manager.webServiceRepository.GetAllWebServicesWithTests(rsdb.GetConnection())
	if err != nil {
		return errors.WithStack(err)
	}

	rslog.Debugf("webServices='%+v'", webServices)
	for _, webService := range webServices {
		webServiceScheduler, err := NewWebServiceScheduler(&webService)
		if err != nil {
			return errors.WithStack(err)
		}
		manager.webServiceSchedulers[webService.Id] = webServiceScheduler
	}
	rslog.Debugf("scheduler.webServices='%+v'", manager.webServiceSchedulers)

	return nil
}

func (manager *webServiceScheduleManager) Close() error {
	rslog.Debug("Closing WebServiceManager...")
	for _, s := range manager.webServiceSchedulers {
		if err := s.Close(); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func NewWebServiceScheduleManager(webServiceRepository repositories.WebServiceRepository, webServiceTestRepository repositories.WebServiceTestRepository) (WebServiceScheduleManager, error) {
	if rsvalid.IsZero(webServiceRepository, webServiceTestRepository) {
		return nil, errors.Wrap(rserrors.ErrInvalidParameter, "WebServiceScheduler")
	}
	return &webServiceScheduleManager{
		webServiceSchedulers:           make(map[interface{}]WebServiceScheduler),
		webServiceRepository:           webServiceRepository,
		webServiceTestResultRepository: webServiceTestRepository,
	}, nil
}

type webServiceScheduler struct {
	webService *models.WebService
	closeChan  chan bool
}

func (schedule *webServiceScheduler) Run() error {
	ticker := schedule.webService.Schedule.GetTicker()
	rslog.Debug("Running...")
	for {
		select {
		case <-ticker.C:
			for _, test := range schedule.webService.Tests {
				res, err := test.Execute(schedule.webService)
				if err != nil {
					return err
				}
				rslog.Debugf("executed test:: id='%v'", test.Id)

				if !test.Assertion.Assert(res) {

				}
			}
		case <-schedule.closeChan:
			rslog.Debugf("webService close:: \tid='%v'", schedule.webService.Id)
			return nil
		}
	}
}

func (schedule *webServiceScheduler) Close() error {
	schedule.closeChan <- true
	close(schedule.closeChan)
	return nil
}

func NewWebServiceScheduler(webService *models.WebService) (WebServiceScheduler, error) {
	if rsvalid.IsZero(webService) {
		return nil, rserrors.ErrInvalidParameter
	}
	return &webServiceScheduler{
		webService: webService,
		closeChan:  make(chan bool, 1),
	}, nil
}
