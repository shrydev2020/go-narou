package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"narou/adapter/logger"
	"narou/adapter/query/hameln"
	"narou/adapter/query/narou"
	"narou/adapter/repository/epub"
	metadataRepo "narou/adapter/repository/metadata"
	"narou/adapter/repository/novel"
	"narou/config"
	"narou/domain/metadata"
	"narou/domain/query"
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
	cfg := config.GetConfigure()
	a := download.NewDownloadInteractor(
		ctx,
		lg,
		metadataRepo.NewRepository(db),
		novel.NewRepository(storage.NewManager(cfg)),
		nil,
		crawl.NewCrawler(lg),
		DependencyInjection(args[0]))

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
		novel.NewRepository(storage.NewManager(cfg)),
		epub.NewRepository(storage.NewManager(cfg)),
		config.GetConfigure(),
		DependencyInjection(args[0]))

	if err := controller.NewConvertController(cvt, lg).Execute(args); err != nil {
		lg.Error("err", "error", err.Error())
	}
	m, _ := metadataRepo.NewRepository(db).FindByTopURI(metadata.URI(args[0]))
	msg := fmt.Sprintf("convert completed. epub generated at %s/%s/%s/%s.epub", storage.NewManager(cfg).GetDist(),
		m.SiteName, m.Title, m.Title)
	lg.Info("finis!", "msg", msg)
	return nil
}

func DependencyInjection(uri string) func(string) (query.IQuery, error) {
	if strings.Contains(uri, "https://ncode.syosetu.com/") {
		return narou.New
	} else if strings.Contains(uri, "https://syosetu.org/") {
		return hameln.New
	}
	panic("not implement")
}
