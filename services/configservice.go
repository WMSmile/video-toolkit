package services

import (
	appconfig "video-toolkit/config"
)

type ConfigService struct {
}

func (c *ConfigService) GetConfig() appconfig.AppConfig {
	return appconfig.AppCfg
}

func (c *ConfigService) SetConfig(config appconfig.AppConfig) {
	appconfig.AppCfg = config
	appconfig.SaveConfig()
}
