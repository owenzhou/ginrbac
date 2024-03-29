package yaml

type Mysql struct {
	Host            string `mapstructure:"host" json:"host" yaml:"host"`
	Port            string `mapstructure:"port" json:"port" yaml:"port"`
	Username        string `mapstructure:"username" json:"username" yaml:"username"`
	Password        string `mapstructure:"password" json:"password" yaml:"password"`
	Dbname          string `mapstructure:"dbname" json:"dbname" yaml:"dbname"`
	Charset         string `mapstructure:"charset" json:"charset" yaml:"charset"`
	MaxOpenConns    int    `mapstructure:"maxOpenConns" json:"maxOpenConns" yaml:"maxOpenConns"`
	MaxIdleConns    int    `mapstructure:"maxIdleConns" json:"maxIdleConns" yaml:"maxIdleConns"`
	ConnMaxLifeTime int    `mapstructure:"connMaxLifeTime" json:"connMaxLifeTime" yaml:"connMaxLifeTime"`
}
