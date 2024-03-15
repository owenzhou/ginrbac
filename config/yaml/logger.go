package yaml

type Logger struct {
	Level      string `mapstructure:"level" json:"level" yaml:"level"`
	Encoding   string `mapstructure:"encoding" json:"encoding" yaml:"encoding"`
	OutputDir  string `mapstructure:"output-dir" json:"output-dir" yaml:"output-dir"`
	MaxAge     int    `mapstructure:"maxage" json:"maxage" yaml:"maxage"`
	MaxSize    int    `mapstructure:"maxsize" json:"maxsize" yaml:"maxsize"`
	MaxBackups int    `mapstructure:"maxbackups" json:"maxbackups" yaml:"maxbackups"`
	Compress   bool   `mapstructure:"compress" json:"compress" yaml:"compress"`
	Filename   string `mapstructure:"filename" json:"filename" yaml:"filename"`
}
