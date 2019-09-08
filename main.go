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
	endpointRepository := repositories.NewEndpointRepository()

	if err := rsdb.CreateTables(webserviceRepository, endpointRepository); err != nil {
		logrus.Fatal(err)
	}

	webServiceService, err := services.NewWebServiceService(endpointRepository)
	if err != nil {
		logrus.Fatal(err)
	}

	endpointService, err := services.NewEndpointService(endpointRepository)
	if err != nil {
		logrus.Fatal(err)
	}

	webServiceHandler, err := handlers.NewWebServiceHandler(webServiceService)
	if err != nil {
		logrus.Fatal(err)
	}

	endpointHandler, err := handlers.NewEndpointHandler(webServiceService, endpointService)
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
				v1WebService.GET("", webServiceHandler.GetWebServiceById)
				v1WebService.DELETE("", webServiceHandler.DeleteWebServiceById)
				v1WebService.PUT("", webServiceHandler.UpdateWebServiceById)
			}

			v1OneWebServiceEndpoints := v1OneWebService.Group("/endpoints")
			{
				v1OneWebServiceEndpoints.POST("", endpointHandler.CreateEndpoint)
			}

		}

		v1Endpoint := v1.Group("/endpoints")
		{
			v1OneEndpoint := v1Endpoint.Group(fmt.Sprintf("/:%s", handlers.EndpointIdParam))
			{
				v1OneEndpoint.GET("", endpointHandler.GetEndpoint)
			}
		}
	}

	e.Logger.Fatal(e.Start(":1323"))
}
