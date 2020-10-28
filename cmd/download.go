package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"narou/adapter/logger"
	"narou/adapter/repository/epub"
	metadata2 "narou/adapter/repository/metadata"
	"narou/adapter/repository/novel"
	"narou/config"
	"narou/domain/metadata"
	"narou/infrastructure/database"
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
		Long:  `start download and convert`,
		Args:  executeArgs,
		Run:   executeDownload,
	}

	rootCmd.AddCommand(downloadCmd)
}

func executeArgs(_ *cobra.Command, args []string) error {
	return validate(args)
}

func executeDownload(_ *cobra.Command, args []string) {
	lg := logger.NewLogger()
	c := config.InitConfigure()
	db, err := database.GetConn()

	if err != nil {
		panic(err)
	}

	a := download.NewDownloadInteractor(
		lg,
		metadata2.NewRepository(db),
		novel.NewRepository(c.GetStorageConfig()),
		nil,
		crawl.NewCrawler(lg))

	ret, err := controller.NewDownloadController(a, lg).Execute(args)
	if err != nil {
		fmt.Printf("error occurred:%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("download completed. episodes:%d", len(ret))
	dist, subDist := c.GetStorageConfig()
	ctx := context.Background()
	cvt := convert.NewConvertInteractor(ctx,
		lg,
		metadata2.NewRepository(db),
		novel.NewRepository(dist, subDist),
		epub.NewRepository(dist),
		c)
	controller.NewConvertController(cvt, lg).Execute(args)
	m, _ := metadata2.NewRepository(db).FindByTopURI(metadata.URI(args[0]))
	fmt.Printf("convert completed. epub generated at %s/%s/%s/%s.epub", dist, m.SiteName, m.Title, m.Title)

}
