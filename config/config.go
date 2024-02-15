package config

import (
	"golang.org/x/oauth2"
)

// AppConfig represents the global configuration object
type AppConfig struct {
	OauthConf          *oauth2.Config
	OauthStateStringGl string
}

// GlobalConfig is the global instance of AppConfig
var GlobalConfig *AppConfig
