package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const grpcAddress = "localhost:18080"

type ServiceContext struct {
	echo.Context
	ServiceConn *grpc.ClientConn
}

// GRPCConnMiddleware conn grpc server
func GRPCConnMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, canceler := context.WithTimeout(c.Request().Context(), time.Second*1)
			defer canceler()

			con, err := grpc.DialContext(ctx, grpcAddress, grpc.WithBlock(), grpc.WithInsecure(), newOption())
			defer wrapper(con.Close)

			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			sc := &ServiceContext{
				Context:     c,
				ServiceConn: con,
			}

			return next(sc)
		}
	}
}
func wrapper(f func() error) {
	if err := f(); err != nil {
		fmt.Println(err)
	}
}

func newOption() grpc.DialOption {
	return grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                GRPCClientKeepaliveTime,
		Timeout:             GRPCClientKeepaliveTimeout,
		PermitWithoutStream: true,
	})
}

const (
	// GRPCClientKeepaliveTime は活動がなくなってから PING を送るまでの間隔を表す。
	GRPCClientKeepaliveTime = 1 * time.Second
	// GRPCClientKeepaliveTimeout は PING 応答を待つ時間を表す。
	GRPCClientKeepaliveTimeout = 5 * time.Second
)
