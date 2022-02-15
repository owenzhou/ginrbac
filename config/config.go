package config

import (
	"github.com/spf13/viper"
	"github.com/owenzhou/ginrbac/contracts"
	"os"
)

func newConfig(app contracts.IApplication) *viper.Viper {
	configPath := "config.yaml"
	if _, err := os.Stat(configPath); err != nil{
		return nil
	}

	viper := app.Make("viper").(*viper.Viper)
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	return viper
}
