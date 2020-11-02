package main

import (
	"github.com/labstack/echo/v4"

	"narou/infrastructure/waf/grpc_client/middleware"
	"narou/infrastructure/waf/grpc_client/novel"
)

func main() {
	e := echo.New()
	e.Use(middleware.GRPCConnMiddleware())
	e.GET("/hello", novel.GetList)

	if err := e.Start(":1324"); err != nil {
		e.Logger.Info("shutting down the server")
	}

}
