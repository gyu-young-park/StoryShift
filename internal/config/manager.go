package config

import (
	"flag"
	"path/filepath"
)

var Manager = newConfigManager(injectConfigParser())

func injectConfigParser() ConfigParser {
	configPath := flag.String("config", "", "path to config file")
	extension := filepath.Ext(*configPath)
	parserMapper := map[string]ConfigParser{
		"yaml": newYamlParser(*configPath),
		"env":  newEnvParser(),
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
