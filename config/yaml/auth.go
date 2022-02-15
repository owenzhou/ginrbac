package yaml

type Auth struct {
	Defaults  map[string]string            `mapstructure:"defaults" json:"defaults" yaml:"defaults"`
	Guards    map[string]map[string]string `mapstructure:"guards" json:"guards" yaml:"guards"`
	Providers map[string]map[string]string `mapstructure:"providers" json:"providers" yaml:"providers"`
}
