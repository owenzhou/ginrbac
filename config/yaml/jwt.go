package yaml

type JWT struct {
	SignKey     string `mapstructure:"sign-key" json:"sign-key" yaml:"sign-key"`
	ExpiresTime int64  `mapstructure:"expires-time" json:"expires-time" yaml:"expires-time"`
}
