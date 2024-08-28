package novel

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"narou/domain/metadata"
	"narou/infrastructure/database"
	"narou/sdk/logger"
	metadata3 "narou/usecase/metadata"
)

func Get(c echo.Context) error {
	lg, err := logger.NewServerLogger(c.Request().Context())
	if err != nil {
		return err
	}
	lg.Info("main start")
	defer lg.Info("main end")

	con, _ := database.GetConn()
	lst, _ := metadata3.NewMetaDataListUseCase(lg, metadata.NewRepository(con), nil).Execute(c.Request().Context())

	return c.JSON(http.StatusOK, lst)
}
