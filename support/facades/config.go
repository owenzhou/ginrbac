package facades

import (
	"fmt"

	"github.com/owenzhou/ginrbac/config/yaml"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Config *yaml.Config

type ConfigFacade struct {
	*Facade
}

func (config *ConfigFacade) GetFacadeAccessor() {
	viper := config.App.Make("config").(*viper.Viper)
	if viper == nil {
		fmt.Println("App make config err: config is nil.")
		return
	}
	if err := viper.ReadInConfig(); err != nil {
		panic("Read config error.")
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
