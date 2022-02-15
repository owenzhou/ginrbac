package yaml

import "time"

type Logger struct {
	Level        string        `mapstructure:"level" json:"level" yaml:"level"`
	Debug        bool          `mapstructure:"debug" json:"debug" yaml:"debug"`
	LinkName     string        `mapstructure:"link-name" json:"link-name" yaml:"link-name"`
	Encoding     string        `mapstructure:"encoding" json:"encoding" yaml:"encoding"`
	OutputDir    string        `mapstructure:"output-dir" json:"output-dir" yaml:"output-dir"`
	MaxAge       time.Duration `mapstructure:"max-age" json:"max-age" yaml:"max-age"`
	RotationTime time.Duration `mapstructure:"rotation-time" json:"rotation-time" yaml:"rotation-time"`
}
