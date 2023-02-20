package yaml

type JWT struct {
	SignKey  string `mapstructure:"sign-key" json:"sign-key" yaml:"sign-key"`
	LifeTime int64  `mapstructure:"lifetime" json:"lifetime" yaml:"lifetime"`
}
