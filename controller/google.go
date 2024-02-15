package controller

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/rammyblog/spendwise/config"
	"github.com/rammyblog/spendwise/services"
)

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	config := config.GlobalConfig
	URL, err := url.Parse(config.OauthConf.Endpoint.AuthURL)
	if err != nil {
		log.Fatal("Parse: " + err.Error())
	}
	parameters := url.Values{}
	parameters.Add("client_id", config.OauthConf.ClientID)
	parameters.Add("scope", strings.Join(config.OauthConf.Scopes, " "))
	parameters.Add("redirect_uri", config.OauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", config.OauthStateStringGl)
	URL.RawQuery = parameters.Encode()
	url := URL.String()
	log.Print(url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func CallBackFromGoogle(w http.ResponseWriter, r *http.Request) {
	config := config.GlobalConfig

	state := r.FormValue("state")

	if state != config.OauthStateStringGl {
		log.Printf("invalid oauth state, expected " + config.OauthStateStringGl + ", got " + state + "\n")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")

	if code == "" {
		log.Fatal("Code not found..")
		w.Write([]byte("Code Not Found to provide AccessToken..\n"))
		reason := r.FormValue("error_reason")
		if reason == "user_denied" {
			w.Write([]byte("User has denied Permission.."))
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}

	} else {
		ctx := context.Background()
		token, err := config.OauthConf.Exchange(ctx, code)
		if err != nil {
			log.Fatal("oauthConfGl.Exchange() failed with " + err.Error() + "\n")
			return
		}

		var response map[string]interface{}

		// Get the user details
		err = services.GetResource(ctx, "https://www.googleapis.com/oauth2/v3/userinfo", &response, token.AccessToken)

		if err != nil {
			log.Fatal("Get: " + err.Error() + "\n")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		log.Printf("User Info>> %+v", response)

		w.Write([]byte("Hello, I'm protected\n"))
		return
	}
}
