package config

import (
	"fmt"
	"os"
	"path/filepath"
)

var Manager = newConfigManager(injectConfigParser())

func injectConfigParser() ConfigParser {
	configPath := os.Getenv("STORY_SHIFT_CONFIG_FILE")
	extension := filepath.Ext(configPath)
	fmt.Println("Config file path: " + configPath)

	parserMapper := map[string]ConfigParser{
		".yaml": newYamlParser(configPath),
		".env":  newEnvParser(),
	}

	parser, ok := parserMapper[extension]
	if !ok {
		return newEnvParser()
	}

	return parser
}

func newConfigManager(configParser ConfigParser) configManager {
	manager := configManager{}
	configParser.parse(&manager.ConfigModel)
	return manager
}

type configManager struct {
	ConfigModel
}
