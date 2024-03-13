package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rammyblog/spendwise/controller"
	localMiddleware "github.com/rammyblog/spendwise/middleware"
	"github.com/rammyblog/spendwise/templates"
	"github.com/rammyblog/spendwise/utils"
)

func router() *chi.Mux {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.RedirectSlashes)

	fs := http.FileServer(http.Dir("static/"))
	r.Handle("/*", http.StripPrefix("/", fs))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		templates.Render(w, "index.html", nil, true)
	})
	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		w.WriteHeader(http.StatusOK)
		errorMsg, _ := utils.GetCookie(r, "errorSw")
		var data struct {
			Error string
		}

		data.Error = errorMsg
		templates.Render(w, "login.html", data, true)
	})
	r.Get("/signup", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		errorMsg, _ := utils.GetCookie(r, "errorSw")

		var data struct {
			Error string
		}

		data.Error = errorMsg

		templates.Render(w, "signup.html", data, true)
	})
	r.Post("/handle-login", controller.HandleGoogleLogin)
	r.Get("/callback-google", controller.CallBackFromGoogle)

	r.Group(func(r chi.Router) {
		r.Use(localMiddleware.IsGoogleAuthenticated)
		r.Get("/dashboard", func(w http.ResponseWriter, r *http.Request) {
			controller.Dashboard(w, r, 5)
		})

		r.Get("/dashboard/add-expense", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Println(r.Referer(), "referer")
			w.WriteHeader(http.StatusOK)
			data := map[string]interface{}{
				"Message": "",
				"Link":    "",
				"Error":   "",
				"Header":  "",
			}
			categories, err := controller.GetCategories()
			if err != nil {
				data["Error"] = "Error getting categories"
				log.Println("Error getting categories: ", err)
				templates.Render(w, "add-expense.html", data, false)
				return
			}
			data["Categories"] = categories
			data["Link"] = "/dashboard/add-expense"
			if strings.Contains(r.Referer(), "expenses") {
				data["Header"] = `{"Expenses-Page": "true"}`
				data["Target"] = "#expenses"

			}
			templates.Render(w, "add-expense.html", data, false)
		})

		r.Post("/dashboard/add-expense", controller.AddExpense)

		r.Get("/dashboard/expense-graph", controller.ExpenseGraph)
		r.Get("/dashboard/expenses", controller.ExpenseList)
		r.Get("/dashboard/expenses/{id}", controller.ExpenseDetail)

	})

	return r
}

// hx-headers='{"Expenses-Page": "true"}'
