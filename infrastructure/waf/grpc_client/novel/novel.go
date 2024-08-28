package novel

import (
	"context"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"

	"narou/infrastructure/waf/grpc_client/middleware"
	"narou/sdk/logger"
	pb "narou/usecase/port/boudary/proto/novel"
)

func GetList(c echo.Context) error {
	ctx, canceler := context.WithTimeout(c.Request().Context(), time.Second*1)
	defer canceler()

	lg, err := logger.NewServerLogger(ctx)
	if err != nil {
		return errors.Wrapf(err, "logger err")
	}
	lg.Info("grpc grpc_client get start")
	lg.Info("ctx", "ctx", ctx)
	defer lg.Info("grpc grpc_client get end")

	sc, _ := c.(*middleware.ServiceContext)

	novel, err := pb.NewNovelListClient(sc.ServiceConn).Get(ctx, &pb.Req{})
	if err != nil {
		panic(err)
	}
	return c.JSON(http.StatusOK, novel)
}
