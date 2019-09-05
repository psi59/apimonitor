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

	webserviceRepository := repositories.NewWebServiceRepository()
	endpointRepository := repositories.NewWebServiceRepository()

	if err := rsdb.CreateTables(webserviceRepository, endpointRepository); err != nil {
		logrus.Fatal(err)
	}

	webServiceService, err := services.NewWebServiceService(endpointRepository)
	if err != nil {
		logrus.Fatal(err)
	}

	webServiceHandler, err := handlers.NewWebServiceHandler(webServiceService)
	if err != nil {
		logrus.Fatal(err)
	}

	v1 := e.Group("/v1")
	{
		v1WebService := v1.Group("/webservices")
		{
			v1WebService.POST("", webServiceHandler.CreateWebService)
			v1WebService.GET("", webServiceHandler.GetWebServiceList)
			v1WebService.GET(fmt.Sprintf("/:%s", handlers.WebServiceIdParam), webServiceHandler.GetWebServiceById)
			v1WebService.DELETE(fmt.Sprintf("/:%s", handlers.WebServiceIdParam), webServiceHandler.DeleteWebServiceById)
			v1WebService.PUT(fmt.Sprintf("/:%s", handlers.WebServiceIdParam), webServiceHandler.UpdateWebServiceById)
		}
	}

	e.Logger.Fatal(e.Start(":1323"))
}
