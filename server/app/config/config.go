package config

import (
	"github.com/jinzhu/configor"
	"path/filepath"
)

var Config = struct {
	App struct {
		SecretKey   string
		BaseDir     string
		LogDir      string
		DataDir     string
		StaticDir   string
		UserDataDir string
	}

	DataBase struct {
		Host     string
		Port     string
		Username string
		Password string
		Name     string
	}

	Docker struct {
		Registry string
	}
}{}

var App = &Config.App
var Database = &Config.DataBase
var Docker = &Config.Docker

func init() {
	if err := configor.Load(&Config, "conf/config.toml"); err != nil {
		panic(err)
	}

	Config.App.BaseDir, _ = filepath.Abs(".")
	//Config.App.DataDir = filepath.Join(Config.App.BaseDir, "data")
	Config.App.LogDir = filepath.Join(Config.App.BaseDir, "log")
	//Config.App.StaticDir = filepath.Join(Config.App.DataDir, "static")
	//Config.App.UserDataDir = filepath.Join(Config.App.DataDir, "user")
}
