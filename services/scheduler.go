package services

import (
	"time"

	"github.com/pkg/errors"
	"github.com/realsangil/apimonitor/pkg/rsmodels"
	"github.com/realsangil/apimonitor/pkg/rsstr"

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

type ScheduleConstructor interface {
	Init() error
}

type Scheduler interface {
	ScheduleRunner
	ScheduleCloser
}

type ScheduleManager interface {
	ScheduleRunner
	ScheduleConstructor
	ScheduleCloser
}

type TestScheduleManager struct {
	testSchedulers       map[string]Scheduler
	testRepository       repositories.TestRepository
	testResultRepository repositories.TestResultRepository
	resultChan           chan *models.TestResult
	closeChan            chan bool
}

func (manager *TestScheduleManager) Run() error {
	rslog.Debug("Running WebServiceManager...")
	errChan := make(chan error, 100)
	for _, s := range manager.testSchedulers {
		go func(s Scheduler, errChan chan<- error) {
			if err := s.Run(); err != nil {
				errChan <- err
			}
		}(s, errChan)
	}
	for {
		select {
		case err := <-errChan:
			rslog.Errorf("error='%v'", err)
		// 	TODO: 에러 프린팅
		case result := <-manager.resultChan:
			rslog.Debugf("result='%+v'", result)
			if err := manager.testResultRepository.Create(rsdb.GetConnection(), result); err != nil {
				errChan <- err
			}
		case <-manager.closeChan:
			rslog.Debug("Closed TestScheduleManager")
			_ = manager.Close()
			return nil
		}
	}
}

func (manager *TestScheduleManager) Init() error {
	if err := manager.Close(); err != nil {
		return errors.WithStack(err)
	}
	if err := manager.initTests(); err != nil {
		return errors.WithStack(err)
	}
	if err := manager.Run(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (manager *TestScheduleManager) initTests() error {
	rslog.Debug("Initialing Scheduler...")
	tests := make([]*models.Test, 0)
	filter := rsdb.ListFilter{
		Page:       -1,
		NumItem:    -1,
		Conditions: nil,
	}
	totalCount, err := manager.testRepository.GetList(rsdb.GetConnection(), &tests, filter, nil)
	if err != nil {
		return errors.WithStack(err)
	}
	rslog.Debugf("test total count='%d'", totalCount)

	rslog.Debugf("tests='%+v'", tests)
	for _, test := range tests {
		testScheduler, err := NewTestScheduler(test, manager.resultChan)
		if err != nil {
			return errors.WithStack(err)
		}
		manager.testSchedulers[test.Id] = testScheduler
	}
	rslog.Debugf("scheduler.testServices='%+v'", manager.testSchedulers)

	return nil
}

func (manager *TestScheduleManager) Close() error {
	rslog.Debug("Closing WebServiceManager...")
	for _, s := range manager.testSchedulers {
		_ = s.Close()
	}
	// manager.closeChan <- true
	return nil
}

func NewTestScheduleManager(testRepository repositories.TestRepository, testResultRepository repositories.TestResultRepository) (ScheduleManager, error) {
	if rsvalid.IsZero(testRepository, testResultRepository) {
		return nil, errors.Wrap(rserrors.ErrInvalidParameter, "Scheduler")
	}
	return &TestScheduleManager{
		testSchedulers:       make(map[string]Scheduler),
		testRepository:       testRepository,
		testResultRepository: testResultRepository,
		resultChan:           make(chan *models.TestResult, 1000),
	}, nil
}

type testScheduler struct {
	test       *models.Test
	closeChan  chan bool
	resultChan chan<- *models.TestResult
}

func (schedule *testScheduler) Run() error {
	ticker := schedule.test.Schedule.GetTicker()
	rslog.Debug("Running...")
	for {
		select {
		case <-ticker.C:
			test := schedule.test
			res, err := test.Execute()
			if err != nil {
				return err
			}
			rslog.Debugf("executed test:: id='%v'", test.Id)
			schedule.resultChan <- &models.TestResult{
				DefaultValidateChecker: rsmodels.ValidatedDefaultValidateChecker,
				Id:                     rsstr.NewUUID(),
				TestId:                 test.Id,
				IsSuccess:              test.Assertion.Assert(res),
				StatusCode:             res.StatusCode,
				Response:               res.Body,
				ResponseTime:           res.ResponseTime,
				TestedAt:               time.Now(),
			}
		case <-schedule.closeChan:
			rslog.Debugf("test close:: \tid='%v'", schedule.test.Id)
			return nil
		}
	}
}

func (schedule *testScheduler) Close() error {
	schedule.closeChan <- true
	close(schedule.closeChan)
	return nil
}

func NewTestScheduler(test *models.Test, resultChan chan<- *models.TestResult) (Scheduler, error) {
	if rsvalid.IsZero(test) {
		return nil, rserrors.ErrInvalidParameter
	}
	return &testScheduler{
		test:       test,
		closeChan:  make(chan bool, 1),
		resultChan: resultChan,
	}, nil
}
