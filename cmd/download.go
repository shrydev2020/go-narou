package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"narou/adapter/logger"
	"narou/adapter/repository/epub"
	metadataRepo "narou/adapter/repository/metadata"
	"narou/adapter/repository/novel"
	"narou/config"
	"narou/domain/metadata"
	"narou/infrastructure/database"
	"narou/infrastructure/storage"
	"narou/interface/controller"
	"narou/interface/gateway/crawl"
	"narou/usecase/interactor/convert"
	"narou/usecase/interactor/download"
)

func init() {
	// downloadCmd represents the download command
	var downloadCmd = &cobra.Command{
		Use:   "d",
		Short: "start download and convert",
		Args:  executeArgs,
		RunE:  executeDownload,
	}

	rootCmd.AddCommand(downloadCmd)
}

func executeArgs(_ *cobra.Command, args []string) error {
	return validate(args)
}

func executeDownload(c *cobra.Command, args []string) error {
	ctx := c.Context()
	lg := logger.NewLogger(ctx)
	db, err := database.GetConn()

	if err != nil {
		return err
	}

	a := download.NewDownloadInteractor(
		ctx,
		lg,
		metadataRepo.NewRepository(db),
		novel.NewRepository(storage.NewManager()),
		nil,
		crawl.NewCrawler(lg))

	ret, err := controller.NewDownloadController(a, lg).Execute(args)
	if err != nil {
		lg.Error("error occurred", "err", err.Error())
		return err
	}
	lg.Info("download completed. episodes", "total", len(ret))
	cvt := convert.NewConvertInteractor(
		ctx,
		lg,
		metadataRepo.NewRepository(db),
		novel.NewRepository(storage.NewManager()),
		epub.NewRepository(storage.NewManager()),
		config.InitConfigure())

	if err := controller.NewConvertController(cvt, lg).Execute(args); err != nil {
		lg.Error("err", "error", err.Error())
	}
	m, _ := metadataRepo.NewRepository(db).FindByTopURI(metadata.URI(args[0]))
	msg := fmt.Sprintf("convert completed. epub generated at %s/%s/%s/%s.epub", storage.NewManager().GetDist(),
		m.SiteName, m.Title, m.Title)
	lg.Info("finis!", "msg", msg)
	return nil
}
