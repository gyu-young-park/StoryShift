package config

type AppConfigModel struct {
	Log    LogConfigModel    `json:"log" yaml:"log"`
	Server ServerConfigModel `json:"server" yaml:"server"`
}

type LogConfigModel struct {
	Library string `json:"library" yaml:"library"`
	Level   string `json:"level" yaml:"level"`
}

type ServerConfigModel struct {
	Port string `json:"port" yaml:"port"`
}

type RedisConfigModel struct {
	Test     bool   `json:"test" yaml:"test"`
	Addr     string `json:"addr" yaml:"addr"`
	Password string `json:"password" yaml:"password"`
}

type VelogConfigModel struct {
	Url    string `json:"url" yaml:"url"`
	ApiUrl string `json:"api_url" yaml:"api_url"`
}

type ConfigModel struct {
	AppConfig   AppConfigModel   `json:"app" yaml:"app"`
	RedisConfig RedisConfigModel `json:"redis" yaml:"redis"`
	VelogConfig VelogConfigModel `json:"velog" yaml:"velog"`
}
