package middleware

import (
	"fmt"
	"net/http"

	"github.com/rammyblog/spendwise/utils"
)

func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := utils.GetCookie(r, "swAccess")

		// TODO: need to fix this
		if err != nil || accessToken == "" {
			fmt.Println("No access token")
			http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
			return
		}

		next.ServeHTTP(w, r)
	})
}
