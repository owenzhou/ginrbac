package yaml

type Session struct {
	SecretKey string `mapstructure:"secret-key" json:"secret-key" yaml:"secret-key"`
	LifeTime  int    `mapstructure:"lifetime" json:"lifetime" yaml:"lifetime"`
}
