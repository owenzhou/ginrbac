package viper

import (
	"github.com/spf13/viper"
)

func newViper() *viper.Viper {
	return viper.New()
}
