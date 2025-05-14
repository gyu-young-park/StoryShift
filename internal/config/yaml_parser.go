package config

import (
	"io"
	"os"

	"github.com/stretchr/testify/assert/yaml"
)

func newYamlParser(configFilePath string) yamlParser {
	return yamlParser{
		configFilePath: configFilePath,
	}
}

type yamlParser struct {
	configFilePath string
}

func (y yamlParser) parse(configModel *ConfigModel) {
	f, err := os.Open(y.configFilePath)
	if err != nil {
		panic("failed to open config file")
	}

	data, err := io.ReadAll(f)
	if err != nil {
		panic("failed to read config file")
	}

	err = yaml.Unmarshal(data, configModel)
	if err != nil {
		panic("failed to unmarshal yaml config data")
	}
}
