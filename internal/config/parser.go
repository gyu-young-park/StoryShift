package config

import "os"

type ConfigParser interface {
	parse(*ConfigModel)
}

func newEnvParser() envParser {
	return envParser{}
}

type envParser struct {
}

func (e envParser) parse(configModel *ConfigModel) {
	configModel.VelogConfig.URL = getEnvDataWithDefault("VELOG_URL", "https://v2.velog.io/graphql")
}

func getEnvDataWithDefault(key, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}

	return v
}
