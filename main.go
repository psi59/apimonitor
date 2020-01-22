package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	"github.com/realsangil/apimonitor/config"
	"github.com/realsangil/apimonitor/handlers"
	"github.com/realsangil/apimonitor/middlewares"
	"github.com/realsangil/apimonitor/pkg/rsdb"
	"github.com/realsangil/apimonitor/pkg/rslog"
	"github.com/realsangil/apimonitor/repositories"
	"github.com/realsangil/apimonitor/services"
)

func main() {
	if err := config.Init(config.ConfigFilePath); err != nil {
		logrus.Fatal(err)
	}

	serverConfig := config.GetServerConfig()

	if err := rslog.Init(&serverConfig.Logger); err != nil {
		logrus.Fatal(err)
	}

	if err := rsdb.Init(&serverConfig.DB); err != nil {
		rslog.Fatal(err)
	}

	e := echo.New()
	e.Use(
		middlewares.ReplaceContextMiddleware,
		middleware.Logger(),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     []string{"http://localhost:3001", "http://localhost:3000"},
			AllowCredentials: true,
			AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		}),
	)
	e.HTTPErrorHandler = middlewares.ErrorHandleMiddleware

	webServiceRepository := repositories.NewWebServiceRepository()
	testRepository := repositories.NewTestRepository()
	testResultRepository := repositories.NewTestResultRepository()

	if err := rsdb.CreateTables(
		webServiceRepository,
		testRepository,
		testResultRepository,
	); err != nil {
		rslog.Fatal(err)
	}

	testSchedulerManager, err := services.NewTestScheduleManager(testRepository, testResultRepository)
	if err != nil {
		rslog.Fatal(err)
	}
	go func() {
		if err := testSchedulerManager.Init(); err != nil {
			rslog.Fatal(err)
		}
	}()

	webServiceService, err := services.NewWebServiceService(webServiceRepository)
	if err != nil {
		rslog.Fatal(err)
	}

	testService, err := services.NewTestService(testRepository, testSchedulerManager)
	if err != nil {
		rslog.Fatal(err)
	}

	testResultService, err := services.NewTestResultService(testResultRepository)
	if err != nil {
		rslog.Fatal(err)
	}

	webServiceHandler, err := handlers.NewWebServiceHandler(webServiceService)
	if err != nil {
		rslog.Fatal(err)
	}

	testHandler, err := handlers.NewTestHandler(webServiceService, testService)
	if err != nil {
		rslog.Fatal(err)
	}

	testResultHandler, err := handlers.NewTestResultHandler(testResultService)
	if err != nil {
		rslog.Fatal(err)
	}

	v1 := e.Group("/v1")
	{
		v1WebService := v1.Group("/webservices")
		{
			v1WebService.POST("", webServiceHandler.CreateWebService)
			v1WebService.GET("", webServiceHandler.GetWebServiceList)

			v1OneWebService := v1WebService.Group(fmt.Sprintf("/:%s", handlers.WebServiceIdParam))
			{
				v1OneWebService.GET("", webServiceHandler.GetWebServiceById)
				v1OneWebService.DELETE("", webServiceHandler.DeleteWebServiceById)
				v1OneWebService.PUT("", webServiceHandler.UpdateWebServiceById)
				v1OneWebService.GET("/results", testResultHandler.GetListByWebService)

				v1Test := v1OneWebService.Group("/tests")
				{
					v1Test.POST("", testHandler.CreateTest)
					v1Test.GET("", testHandler.GetTestList)
				}
			}
		}

		v1OneTest := v1.Group(fmt.Sprintf("/tests/:%s", handlers.TestIdParam))
		{
			v1OneTest.GET("", testHandler.GetTest)
			v1OneTest.DELETE("", testHandler.DeleteTest)
			v1OneTest.PUT("", testHandler.UpdateTest)
			v1OneTest.GET("/execute", testHandler.ExecuteTest)
			v1OneTest.GET("/results", testResultHandler.GetListByTest)
		}
	}

	e.Logger.Fatal(e.Start(":1323"))
}
