package config

type Logistics struct {
	URL   string `mapstructure:"url" json:"url" yaml:"url"`
	CpURL string `mapstructure:"cp-url" json:"cp-url" yaml:"cp-url"`
	Token string `mapstructure:"token" json:"token" yaml:"token"`
}
