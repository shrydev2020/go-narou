package cmd

import (
	"fmt"
	"os"

	"narou/config"
	"narou/infrastructure/database"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "narou",
	Short: "scraping from narou and so on...",
	Long:  `scraping from narou、and generate epub`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	c := config.InitConfigure()
	if err := database.OpenDB(c.GetDBConfig()); err != nil {
		fmt.Printf("error occurred:%s", err.Error())
		os.Exit(1)
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
