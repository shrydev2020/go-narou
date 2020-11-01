package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"narou/infrastructure/waf/server/novel"
)

type server struct{}

func New() *server {
	return &server{}
}

// todo use server logger, server group
func (server) Start() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/novel", novel.GetList)

	e.Logger.Fatal(e.Start(":1323"))
}
