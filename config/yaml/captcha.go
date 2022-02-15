package yaml

type Captcha struct {
	Width     int    `mapstructure:"width" json:"width" yaml:"width"`
	Height    int    `mapstructure:"height" json:"height" yaml:"height"`
	Num       int    `mapstructure:"num" json:"num" yaml:"num"`
	SecretKey string `mapstructure:"secret-key" json:"secret-key" yaml:"secret-key"`
	Expired   int64  `mapstructure:"expired" json:"expired" yaml:"expired"`
}
