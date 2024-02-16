package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rammyblog/spendwise/config"
	"github.com/rammyblog/spendwise/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	port := fmt.Sprintf(":%v", os.Getenv("PORT"))

	config.GlobalConfig = &config.AppConfig{}

	// Set google variables
	services.InitializeOAuthGoogle()

	handler := router()

	server := &http.Server{
		Addr:    port,
		Handler: handler,
	}
	log.Printf("[info] start http server listening %s", port)
	server.ListenAndServe()

}
