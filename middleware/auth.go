package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rammyblog/spendwise/utils"
)

type contextKey struct {
	name string
}

var UserIDKey = &contextKey{"userID"}

func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := utils.GetCookie(r, "swAccess")
		// TODO: need to fix this
		if err != nil || accessToken == "" {
			fmt.Println(err)
			http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetUserIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := utils.GetCookie(r, "usw")
		if err != nil {
			HandleError(w, err, "Error getting user id")
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
