package controller

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rammyblog/spendwise/config"
	"github.com/rammyblog/spendwise/models"
	"github.com/rammyblog/spendwise/repositories"
	"github.com/rammyblog/spendwise/services"
	"github.com/rammyblog/spendwise/utils"
)

type UserInfo struct {
	Sub           string    `json:"sub"`
	GivenName     string    `json:"given_name"`
	FamilyName    string    `json:"family_name"`
	Nickname      string    `json:"nickname"`
	Name          string    `json:"name"`
	Picture       string    `json:"picture"`
	Locale        string    `json:"locale"`
	UpdatedAt     time.Time `json:"updated_at"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
}

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	config := config.GlobalConfig
	URL, err := url.Parse(config.OauthConf.Endpoint.AuthURL)
	if err != nil {
		log.Print("Parse: " + err.Error())
	}
	parameters := url.Values{}
	parameters.Add("client_id", config.OauthConf.ClientID)
	parameters.Add("scope", strings.Join(config.OauthConf.Scopes, " "))
	parameters.Add("redirect_uri", config.OauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	referer := r.Header.Get("Referer")
	parameters.Add("state", "referer="+referer+"?state="+config.OauthStateStringGl)
	URL.RawQuery = parameters.Encode()
	url := URL.String()
	w.Header().Set("HX-Redirect", url)

}

func CallBackFromGoogle(w http.ResponseWriter, r *http.Request) {
	config := config.GlobalConfig

	state := r.FormValue("state")

	if !strings.Contains(state, config.OauthStateStringGl) {
		log.Printf("invalid oauth state, expected " + config.OauthStateStringGl + ", got " + state + "\n")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")

	if code == "" {
		log.Print("Code not found..")
		reason := r.FormValue("error_reason")
		if reason == "user_denied" {
			w.Write([]byte("User has denied Permission.."))
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}
	} else {
		ctx := context.Background()
		token, err := config.OauthConf.Exchange(ctx, code)
		if err != nil {
			log.Print("oauthConfGl.Exchange() failed with " + err.Error() + "\n")
			return
		}

		var response UserInfo

		// Get the user details
		err = services.GetResource(ctx, "https://www.googleapis.com/oauth2/v3/userinfo", &response, token.AccessToken)

		if err != nil {
			log.Print("Get: " + err.Error() + "\n")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		// Get the user repository

		isLogin := strings.Contains(state, "login")

		isRegister := strings.Contains(state, "signup")

		if isLogin {
			userRepo := repositories.NewUserRepository(config.DB)
			user, err := userRepo.FindByEmail(response.Email)
			if err != nil {
				log.Print("User not found: " + err.Error() + "\n")
				utils.SetCookie(w, "errorSw", "user not found", time.Now().Add(10*time.Second))
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			}

			// write access token and refresh token to cookie
			utils.SetCookie(w, "usw", user.ID.String(), token.Expiry)
			utils.SetCookie(w, "access_token", token.AccessToken, token.Expiry)
			utils.SetCookie(w, "refresh_token", token.RefreshToken, token.Expiry)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

		}

		if isRegister {
			userRepo := repositories.NewUserRepository(config.DB)
			user := &models.User{
				Email:         response.Email,
				FirstName:     response.GivenName,
				LastName:      response.FamilyName,
				Picture:       response.Picture,
				Provider:      "google",
				EmailVerified: response.EmailVerified,
				ProviderID:    response.Sub,
			}
			err = userRepo.Create(user)

			if err != nil {
				if strings.Contains(err.Error(), `duplicate key value violates unique constraint "users_email_key"`) {
					// user already exists
					utils.SetCookie(w, "errorSw", "user already exist, kindly login", time.Now().Add(20*time.Second))
					http.Redirect(w, r, "/signup", http.StatusTemporaryRedirect)
					return
				}
				log.Print("Error creating user: " + err.Error() + "\n")
				utils.SetCookie(w, "errorSw", "Error creating user", time.Now().Add(10*time.Second))
				http.Redirect(w, r, "/signup", http.StatusTemporaryRedirect)
			}
		}

		// redirect to the home page
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

	}
}
