package middleware

import (
	"fmt"
	"log"
	"net/http"
)

func HandleError(w http.ResponseWriter, err error, errorMessage string) {
	log.Println("Error: ", err)
	w.WriteHeader(http.StatusInternalServerError)
	errorMessage = fmt.Sprintf(`{"error": "%s"}`, errorMessage)
	http.Error(w, errorMessage, http.StatusInternalServerError)
}
