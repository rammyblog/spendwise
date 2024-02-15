package services

import (
	"os"

	"github.com/rammyblog/spendwise/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

/*
InitializeOAuthGoogle Function
*/
func InitializeOAuthGoogle() {
	config.GlobalConfig.OauthConf = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8081/callback-google",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	config.GlobalConfig.OauthStateStringGl = os.Getenv("STATE_STRING_GL")
}
