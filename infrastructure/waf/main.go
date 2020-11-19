package main

import (
	"fmt"
	"os"

	"narou/config"
	"narou/infrastructure/database"
	"narou/infrastructure/waf/server"
)

func main() {
	var cfg = config.GetConfigure()
	if err := database.OpenDB(cfg.GetDBConfig()); err != nil {
		fmt.Printf("error occurred:%s", err.Error())
		os.Exit(1)
	}
	server.New().Start()
}
