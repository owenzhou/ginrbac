package yaml

type Config struct {
	Debug    bool                `mapstructure:"debug" json:"debug" yaml:"debug"`
	AppName  string              `mapstructure:"appName" json:"appName" yaml:"appName"`
	Template []map[string]string `mapstructure:"template" json:"template" yaml:"template"`
	Auth     Auth                `mapstructure:"auth" json:"auth" yaml:"auth"`
	Logger   Logger              `mapstructure:"logger" json:"logger" yaml:"logger"`
	Captcha  Captcha             `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
	Mysql    Mysql               `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	JWT      JWT                 `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Redis    Redis               `mapstructure:"redis" json:"redis" yaml:"redis"`
}
