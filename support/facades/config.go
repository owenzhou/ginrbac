package facades

import (
	"fmt"
	"ginrbac/config/yaml"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Config *yaml.Config

type ConfigFacade struct {
	*Facade
}

func (config *ConfigFacade) GetFacadeAccessor() {
	viper := config.App.Make("config").(*viper.Viper)
	if err := viper.ReadInConfig(); err != nil {
		panic("read config error")
	}

	viper.WatchConfig()

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		if err := viper.Unmarshal(&Config); err != nil {
			panic(err)
		}
	})

	if err := viper.Unmarshal(&Config); err != nil {
		panic(err)
	}
}
