package main

import (
	"fmt"

	"github.com/labstack/echo"
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

	txMiddleware := middlewares.NewTransactionMiddleware()

	v1 := e.Group("/v1")
	{
		v1WebService := v1.Group("/webservices")
		{
			v1WebService.POST("", webServiceHandler.CreateWebService, txMiddleware.Tx)
			v1WebService.GET(fmt.Sprintf("/:%s", handlers.WebServiceIdParam), webServiceHandler.GetWebServiceById, txMiddleware.Tx)
		}
	}

	e.Logger.Fatal(e.Start(":1323"))
}
