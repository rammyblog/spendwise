package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rammyblog/spendwise/config"
	"github.com/rammyblog/spendwise/utils"
	"golang.org/x/oauth2"
)

type contextKey struct {
	name string
}

var UserIDKey = &contextKey{"userID"}

func IsGoogleAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config := config.GlobalConfig

		// Get the OAuth token from the request cookie.
		cookie, err := utils.GetCookie(r, "swAccess")
		fmt.Println(err, "err")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		_, err = config.OauthConf.TokenSource(context.Background(), &oauth2.Token{
			AccessToken:  cookie,
			TokenType:    "Bearer",
			RefreshToken: "",
		}).Token()

		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		userId, err := utils.GetCookie(r, "usw")
		if err != nil {
			HandleError(w, err, "Error getting user id")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
