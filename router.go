package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rammyblog/spendwise/controller"
	"github.com/rammyblog/spendwise/templates"
	"github.com/rammyblog/spendwise/utils"
)

type TemplateData struct {
	Error string
	Data  interface{}
}

func router() *chi.Mux {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		templates.Render(w, "index.html", nil)
	})
	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		errorMsg, err := utils.GetAndDeleteCookie(w, r, "errorSw")
		if err != nil {
			log.Println("Error getting cookie: ", err)
		}
		var data struct {
			Error string
		}

		data.Error = errorMsg
		templates.Render(w, "login.html", data)
	})
	r.Get("/signup", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		errorMsg, err := utils.GetAndDeleteCookie(w, r, "errorSw")

		if err != nil {
			log.Println("Error getting cookie: ", err)
		}

		var data struct {
			Error string
		}

		data.Error = errorMsg

		templates.Render(w, "signup.html", data)
	})
	r.Post("/handle-login", controller.HandleGoogleLogin)
	r.Get("/callback-google", controller.CallBackFromGoogle)

	r.Get("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		templates.Render(w, "dashboard.html", nil)
	})

	r.Get("/dashboard/add-expense", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		data := TemplateData{}
		categories, err := controller.GetCategories()
		if err != nil {
			data.Error = "Error getting categories"
			log.Println("Error getting categories: ", err)
			templates.Render(w, "add-expense.html", data)
			return
		}
		data.Data = categories
		errorMsg, err := utils.GetAndDeleteCookie(w, r, "errorSw")
		if err != nil {
			log.Println("Error getting cookie: ", err)
		}
		data.Error = errorMsg
		templates.Render(w, "add-expense.html", data)
	})

	r.Post("/dashboard/add-expense", controller.AddExpense)

	return r
}
