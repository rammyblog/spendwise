package utils

import (
	"net/http"
	"time"
)

func SetCookie(w http.ResponseWriter, name string, value string, expirationTime time.Time) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  expirationTime,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}
