package config

import (
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

// AppConfig represents the global configuration object
type AppConfig struct {
	OauthConf          *oauth2.Config
	OauthStateStringGl string
	DB                 *gorm.DB
}

// GlobalConfig is the global instance of AppConfig
var GlobalConfig *AppConfig
