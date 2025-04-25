package config

type AppConfigModel struct {
	Log LogConfigModel `json:"log"`
}

type LogConfigModel struct {
	Library string `json:"library"`
	Level   string `json:"level"`
}

type VelogConfigModel struct {
	URL string `json:"url"`
}

type ConfigModel struct {
	AppConfig   AppConfigModel   `json:"app"`
	VelogConfig VelogConfigModel `json:"velog"`
}
