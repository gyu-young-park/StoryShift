package config

var Manager = newConfigManager(newEnvParser())

func newConfigManager(configParser ConfigParser) configManager {
	manager := configManager{}
	configParser.parse(&manager.ConfigModel)
	return manager
}

type configManager struct {
	ConfigModel
}

func (c *configManager) Parse(parser ConfigParser) {
	parser.parse(&c.ConfigModel)
}
