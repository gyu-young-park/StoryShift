package config

type VelogConfigModel struct {
	URL string `json:"url"`
}

type ConfigModel struct {
	VelogConfig VelogConfigModel `json:"velog"`
}
