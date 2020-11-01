package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"

	narouMiddleware "narou/infrastructure/waf/server/middleware"
	"narou/infrastructure/waf/server/novel"
)

type server struct{}

func New() *server {
	return &server{}
}

// todo use server logger, server group
func (server) Start() {
	e := echo.New()
	s := grpc.NewServer()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(narouMiddleware.GRPCMiddleware(s))
	e.GET("/novel", novel.GetList)

	e.Logger.Fatal(e.Start(":1323"))
}
