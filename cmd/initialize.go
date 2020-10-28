package cmd

import (
	"bufio"
	"fmt"
	"os"

	"narou/adapter/logger"
	"narou/infrastructure/database"
	"narou/usecase/interactor/initialize"

	"narou/interface/controller"

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
		Run: executeInitialize,
	}

	rootCmd.AddCommand(initCmd)
}

func executeInitialize(c *cobra.Command, _ []string) {
	if Question("Execute DB initialization y/n] ") {
		fmt.Println("execute initialization")
		lg := logger.NewLogger()
		db, _ := database.GetConn()
		i := initialize.NewInitializeInteractor(db)
		if err := controller.NewInitializeController(i, lg, db); err != nil {
			return
		}

		fmt.Println("finish initialization")
	} else {
		fmt.Println("Bye!")
	}
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
