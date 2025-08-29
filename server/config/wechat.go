package config

type Wechat struct {
	AppID          string `mapstructure:"appid" json:"appid" yaml:"appid"`
	Secret         string `mapstructure:"secret" json:"secret" yaml:"secret"`
	SessionUrl     string `mapstructure:"session-url" json:"session-url" yaml:"session-url"`
	AccessTokenUrl string `mapstructure:"access-token-url" json:"access-token-url" yaml:"access-token-url"`
	TelephoneUrl   string `mapstructure:"telephone-url" json:"telephone-url" yaml:"telephone-url"`
}
