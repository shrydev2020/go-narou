package cmd

import (
	"fmt"

	"narou/infrastructure/waf/server"

	"github.com/spf13/cobra"
)

func init() {
	// webCmd represents the web command.
	var webCmd = &cobra.Command{
		Use:   "web",
		Short: "start web server",
		Long:  `start web server and do something..`,
		Run: func(_ *cobra.Command, _ []string) {
			go server.New().Start()
			fmt.Println("web called")
		},
	}
	rootCmd.AddCommand(webCmd)
}
