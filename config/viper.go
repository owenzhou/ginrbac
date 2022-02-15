package config

import (
	"github.com/spf13/viper"
)

func newConfig(config ...string) *viper.Viper {
	var configPath string
	if len(config) == 0 {
		configPath = "config.yaml"
	}
	viper := viper.New()
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	return viper
}
