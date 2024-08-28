package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"narou/config"
	"narou/domain/epub"
	"narou/domain/metadata"
	novel2 "narou/domain/novel"
	"narou/domain/text_query"
	"narou/domain/text_query/hameln"
	"narou/domain/text_query/narou"
	"narou/infrastructure/database"
	"narou/infrastructure/storage"
	"narou/interface/gateway/crawl"
	"narou/sdk/logger"
	"narou/usecase/convert"
	"narou/usecase/download"
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
	lg, err := logger.NewCLILogger(ctx)
	if err != nil {
		return err
	}
	db, err := database.GetConn()
	if err != nil {
		return err
	}
	cfg := config.GetConfigure()
	a := download.NewDownloadUseCase(
		lg,
		metadata.NewRepository(db),
		novel2.NewRepository(storage.NewManager(cfg)),
		crawl.NewCrawler(lg),
		DependencyInjection(args[0]))

	ret, err := a.Execute(metadata.URI(args[0]))
	if err != nil {
		lg.Error("error occurred", "err", err.Error())
		return err
	}
	lg.Info("download completed. episodes", "total", len(ret))
	cvt := convert.NewConvertUseCase(
		lg,
		metadata.NewRepository(db),
		novel2.NewRepository(storage.NewManager(cfg)),
		epub.NewRepository(storage.NewManager(cfg)),
		config.GetConfigure(),
		DependencyInjection(args[0]))

	if err := cvt.Execute(ctx, metadata.URI(args[0])); err != nil {
		lg.Error("err", "error", err.Error())
	}
	m, _ := metadata.NewRepository(db).FindByTopURI(metadata.URI(args[0]))
	msg := fmt.Sprintf("convert completed. epub generated at %s/%s/%s/%s.epub", storage.NewManager(cfg).GetDist(),
		m.SiteName, m.Title, m.Title)
	lg.Info("finis!", "msg", msg)
	return nil
}

func DependencyInjection(uri string) func(string) (text_query.IQuery, error) {
	if strings.Contains(uri, "https://ncode.syosetu.com/") {
		return narou.NewNarouQuery
	} else if strings.Contains(uri, "https://syosetu.org/") {
		return hameln.NewHamelnQuery
	}
	panic("not implement")
}
