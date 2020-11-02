package novel

import (
	"net/http"

	"github.com/labstack/echo/v4"

	metadata3 "narou/usecase/interactor/metadata"

	"narou/adapter/logger"
	metadata2 "narou/adapter/repository/metadata"
	"narou/infrastructure/database"
)

func GetList(c echo.Context) error {
	lg := logger.NewLogger(c.Request().Context())
	lg.Info("main start")
	defer lg.Debug("main end")

	con, _ := database.GetConn()
	lst, _ := metadata3.NewMetaDataListInteractor(c.Request().Context(),
		lg, metadata2.NewRepository(con),
		nil).Execute()

	return c.JSON(http.StatusOK, lst)
}
