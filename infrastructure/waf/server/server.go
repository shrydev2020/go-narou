package server

import (
	"net"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"

	"narou/infrastructure/waf/server/novel"
	pb "narou/usecase/port/boudary/proto/novel"
)

type server struct{}

func New() *server {
	return &server{}
}

// todo use server logger, server group
func (server) Start() {
	e := echo.New()
	grpcServer := grpc.NewServer()
	pb.RegisterNovelListServer(grpcServer, novel.NewGrpcService())
	// TODO
	//go func() {
	//		defer grpcServer.GracefulStop()
	//		<-ctx.Done()
	//	}()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/novel", novel.Get)

	go func() {
		lis, err := net.Listen("tcp", ":18080")
		if err != nil {
			panic(err)
		}
		if err := grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	// Start server
	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("shutting down the server")
	}
}
