package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"narou/domain/metadata"
	"narou/infrastructure/database"
	"narou/sdk/logger"
	"narou/usecase/initialize"

	"github.com/spf13/cobra"
)

func init() {
	// initCmd represents the init command
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "initialize database",
		Long: `
initialize database
`,
		RunE: executeInitialize,
	}

	rootCmd.AddCommand(initCmd)
}

func executeInitialize(c *cobra.Command, _ []string) error {
	if Question("Execute DB initialization y/n] ") {
		fmt.Println("execute initialization")
		lg, err := logger.NewCLILogger(context.Background())
		if err != nil {
			return err
		}
		db, err := database.GetConn()
		if err != nil {
			return err
		}
		if err := initialize.NewInitUseCase(metadata.NewRepository(db)).Execute(); err != nil {
			lg.Error("err", "err", err.Error())
			return err
		}
		fmt.Println("finish initialization")
	} else {
		fmt.Println("Bye!")
	}
	return nil
}

func Question(q string) bool {
	result := true

	fmt.Print(q)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		i := scanner.Text()

		if i == "Y" || i == "y" {
			break
		} else if i == "N" || i == "n" {
			result = false
			break
		} else {
			fmt.Println("enter y or n.")
			fmt.Print(q)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return result
}
