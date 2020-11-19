package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"narou/adapter/logger"
	"narou/adapter/repository/epub"
	metadataRepo "narou/adapter/repository/metadata"
	"narou/adapter/repository/novel"
	"narou/config"
	"narou/infrastructure/database"
	"narou/infrastructure/storage"
	"narou/interface/controller"
	"narou/usecase/interactor/convert"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "convert html to epub",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("need id")
		}
		ctx := c.Context()
		lg := logger.NewLogger(ctx)
		db, err := database.GetConn()

		if err != nil {
			panic(err)
		}

		cvt := convert.NewConvertInteractor(ctx,
			lg,
			metadataRepo.NewRepository(db),
			novel.NewRepository(storage.NewManager()),
			epub.NewRepository(storage.NewManager()),
			config.InitConfigure())

		return controller.NewConvertController(cvt, lg).Execute(args)
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)
}
