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
		logrus.Fatal(err)
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
	webServiceTestRepository := repositories.NewWebServiceTestRepository()

	if err := rsdb.CreateTables(webServiceRepository, webServiceTestRepository); err != nil {
		logrus.Fatal(err)
	}

	webServiceService, err := services.NewWebServiceService(webServiceRepository)
	if err != nil {
		logrus.Fatal(err)
	}

	webServiceTestService, err := services.NewWebServiceTestService(webServiceTestRepository)
	if err != nil {
		logrus.Fatal(err)
	}

	webServiceHandler, err := handlers.NewWebServiceHandler(webServiceService)
	if err != nil {
		logrus.Fatal(err)
	}

	webServiceTestHandler, err := handlers.NewWebServiceTestHandler(webServiceService, webServiceTestService)
	if err != nil {
		logrus.Fatal(err)
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
			}

			v1OneWebServiceWebServiceTests := v1OneWebService.Group("/webServiceTests")
			{
				v1OneWebServiceWebServiceTests.POST("", webServiceTestHandler.CreateWebServiceTest)
			}

		}

		v1WebServiceTest := v1.Group("/tests")
		{
			v1WebServiceTest.GET("", webServiceTestHandler.GetWebServiceTestList)

			v1OneWebServiceTest := v1WebServiceTest.Group(fmt.Sprintf("/:%s", handlers.WebServiceTestIdParam))
			{
				v1OneWebServiceTest.GET("", webServiceTestHandler.GetWebServiceTest)
				v1OneWebServiceTest.DELETE("", webServiceTestHandler.DeleteWebServiceTest)
			}
		}
	}

	e.Logger.Fatal(e.Start(":1323"))
}
