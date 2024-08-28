package cmd

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"

	"narou/config"
	"narou/domain/epub"
	"narou/domain/metadata"
	"narou/domain/novel"
	"narou/infrastructure/database"
	"narou/infrastructure/storage"
	"narou/sdk/logger"
	"narou/usecase/convert"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "convert html to epub",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.Wrapf(ErrRequiredArgsNotFound, "command:%s", "convert")
		}
		ctx := c.Context()
		ctx = context.WithValue(ctx, "batch-name", "convert")
		lg, err := logger.NewCLILogger(ctx)
		if err != nil {
			return err
		}
		db, err := database.GetConn()

		if err != nil {
			return errors.Wrapf(err, "db connection fail. command:%s", "convert")
		}
		cfg := config.GetConfigure()
		useCase := convert.NewConvertUseCase(
			lg,
			metadata.NewRepository(db),
			novel.NewRepository(storage.NewManager(cfg)),
			epub.NewRepository(storage.NewManager(cfg)),
			cfg,
			DependencyInjection(args[0]))
		return useCase.Execute(ctx, metadata.URI(args[0]))
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)
}
