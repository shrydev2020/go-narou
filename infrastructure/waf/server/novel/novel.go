package novel

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"narou/adapter/logger"
	"narou/adapter/repository/epub"
	metadata2 "narou/adapter/repository/metadata"
	"narou/adapter/repository/novel"
	"narou/config"
	"narou/domain/metadata"
	"narou/infrastructure/database"
	"narou/infrastructure/storage"
	"narou/interface/controller"
	"narou/usecase/interactor/convert"
	"narou/usecase/interactor/initialize"
)

func GetList(c echo.Context) error {
	lg := logger.NewLogger(c.Request().Context())
	var cfg = config.InitConfigure()

	lg.Info("main start")
	defer lg.Debug("main end")

	con, _ := database.GetConn()
	i := initialize.NewInitializeInteractor(con)
	b := controller.NewInitializeController(i, lg, con)
	if b != nil {
		panic(b)
	}

	ctx := c.Request().Context()
	it := convert.NewConvertInteractor(
		ctx,
		lg,
		metadata2.NewRepository(con),
		novel.NewRepository(storage.NewManager()),
		epub.NewRepository(storage.NewManager()),
		cfg)

	err := controller.NewConvertController(it, lg).
		Execute([]string{"https://ncode.syosetu.com/n5378gc/"})
	if err != nil {
		lg.Error(err.Error())
	}

	dummy := metadata.Novel{
		ID:              1,
		Author:          "aaa",
		Title:           "",
		FileTitle:       "",
		TopUrl:          "",
		SiteName:        "",
		NovelType:       0,
		End:             false,
		LastUpdate:      time.Time{},
		NewArrivalsDate: time.Time{},
		UseSubdirectory: false,
		GeneralFirstUp:  time.Time{},
		NovelUpdatedAt:  time.Time{},
		GeneralKastUp:   time.Time{},
		Length:          0,
		Suspend:         false,
		GeneralAllNo:    0,
		LastCheckAt:     time.Time{},
	}
	return c.JSON(http.StatusOK, dummy)
}
