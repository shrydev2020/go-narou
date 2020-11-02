package server

import (
	"net"

	"github.com/labstack/echo/v4/middleware"

	narouMiddleware "narou/infrastructure/waf/server/middleware"

	"github.com/labstack/echo/v4"
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

	lis, err := net.Listen("tcp", ":18080")
	if err != nil {
		panic(err)
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(narouMiddleware.GRPCMiddleware(grpcServer))
	e.GET("/novel", novel.GetList)

	go func() {
		if err2 := grpcServer.Serve(lis); err2 != nil {
			panic(err2)
		}
	}()

	// Start server
	if err := e.Start(":1323"); err != nil {
		e.Logger.Info("shutting down the server")
	}
}
