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
		// HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func GetCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func GetAndDeleteCookie(w http.ResponseWriter, r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	DeleteCookie(w, name)
	return cookie.Value, nil
}

func DeleteCookie(w http.ResponseWriter, name string) {
	cookie := http.Cookie{
		Name:  name,
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)
}
