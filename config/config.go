package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

// db setting.
type Database struct {
	Path string
}

// destination to row html.
type Storage struct {
	Dist    string
	SubDist string
}
type Epub struct {
	Lang string
	Ppd  string
}

type config struct {
	Database
	Storage
	Epub
}

func (c *config) GetDBConfig() string {
	return c.Database.Path
}

func (c *config) GetStorageConfig() (string, string) {
	return c.Dist, c.SubDist
}
func (c *config) GetEpubSetting() (string, string) {
	return c.Lang, c.Ppd
}

type IConfigure interface {
	GetDBConfig() string
	GetStorageConfig() (dist, subDist string)
	GetEpubSetting() (lang, ppd string)
}

var (
	once sync.Once
	c    config
)

func GetConfigure() IConfigure {
	initConfigure()
	return &c
}

func initConfigure() {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yml")
		viper.AddConfigPath(".")
		viper.WatchConfig()
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("config file read error")
			fmt.Println(err)
			panic(err)
		}

		if err := viper.Unmarshal(&c); err != nil {
			fmt.Println("config file Unmarshal error")
			fmt.Println(err)
			panic(err)
		}
	})
}
