package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

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
	GetStorageConfig() (dist, subdist string)
	GetEpubSetting() (lang, ppd string)
}

func InitConfigure() IConfigure {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file has changed:", e.Name)
	})

	// 環境変数から設定値を上書きできるように設定
	viper.AutomaticEnv()

	// conf読み取り
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("config file read error")
		fmt.Println(err)
		panic(err)
	}

	// 外部からconfの中身を参照できるようにする
	var c config
	// UnmarshalしてCにマッピング
	if err := viper.Unmarshal(&c); err != nil {
		fmt.Println("config file Unmarshal error")
		fmt.Println(err)
		panic(err)
	}

	return &c
}
