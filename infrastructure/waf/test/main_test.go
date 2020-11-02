package main

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"

	pb "narou/usecase/port/boudary/proto/novel"
)

func main() {
	e := echo.New()
	e.Use(serviceContextMiddleware("localhost:18080"))
	e.GET("/hello", getList)

	if err := e.Start(":1324"); err != nil {
		e.Logger.Info("shutting down the server")
	}

}

func getList(c echo.Context) error {
	sc, _ := c.(*ServiceContext)
	novel, err := sc.ServiceClient.Get(context.Background(), &pb.Req{})
	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, novel)
}

type ServiceContext struct {
	echo.Context
	ServiceClient pb.NovelListClient
}

// このMiddlewareでGRPCにダイアル
func serviceContextMiddleware(grpcAddr string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc, err := grpc.Dial(grpcAddr, grpc.WithBlock(), grpc.WithInsecure())
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			defer cc.Close()

			sc := &ServiceContext{
				Context:       c,
				ServiceClient: pb.NewNovelListClient(cc),
			}

			return next(sc)
		}
	}
}
