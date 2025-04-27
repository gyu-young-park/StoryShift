package config

type AppConfigModel struct {
	Log    LogConfigModel    `json:"log"`
	Server ServerConfigModel `json:"server"`
}

type LogConfigModel struct {
	Library string `json:"library"`
	Level   string `json:"level"`
}

type ServerConfigModel struct {
	Port string `json:"port"`
}

type VelogConfigModel struct {
	URL string `json:"url"`
}

type ConfigModel struct {
	AppConfig   AppConfigModel   `json:"app"`
	VelogConfig VelogConfigModel `json:"velog"`
}
