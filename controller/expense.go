package controller

import (
	"log"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/rammyblog/spendwise/config"
	"github.com/rammyblog/spendwise/models"
	"github.com/rammyblog/spendwise/repositories"
	"github.com/rammyblog/spendwise/templates"
	"github.com/rammyblog/spendwise/utils"
)

func AddExpense(w http.ResponseWriter, r *http.Request) {
	var expense models.Expense

	err := r.ParseForm()
	if err != nil {
		log.Println("Error decoding form: ", err)
		data := map[string]string{
			"Error": "Error adding expense",
		}
		templates.Render(w, "add-expense.html", data)
		return
	}

	decoder := schema.NewDecoder()
	err = decoder.Decode(&expense, r.PostForm)
	if err != nil {
		log.Println("Error decoding form: ", err)
		data := map[string]string{
			"Error": "Error adding expense",
		}
		templates.Render(w, "add-expense.html", data)
		return
	}

	userId, err := utils.GetCookie(r, "usw")

	if err != nil {
		log.Println("Error getting user id: ", err)
		data := map[string]string{
			"Error": "Error adding expense",
		}
		templates.Render(w, "add-expense.html", data)
		return
	}
	expense.UserID = userId

	expenseRepo := repositories.NewExpenseRepository(config.GlobalConfig.DB)
	err = expenseRepo.Create(&expense)
	if err != nil {
		log.Println("Error adding expense: ", err)
		data := map[string]string{
			"Error": "Error adding expense",
		}
		templates.Render(w, "add-expense.html", data)
		return
	}

	data := map[string]string{
		"Message": "Expense added successfully",
		"Link":    "/dashboard/add-expense/1",
	}
	templates.Render(w, "add-expense.html", data)

}
