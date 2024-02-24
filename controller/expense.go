package controller

import (
	"encoding/json"
	"fmt"
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
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, `{"error": "Error adding expense"}`, http.StatusInternalServerError)
		return
	}

	decoder := schema.NewDecoder()
	err = decoder.Decode(&expense, r.PostForm)
	if err != nil {
		log.Println("Error decoding form: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, `{"error": "Error adding expense"}`, http.StatusInternalServerError)
		return
	}

	userId, err := utils.GetCookie(r, "usw")

	if err != nil {
		log.Println("Error getting user id: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, `{"error": "Error adding expense"}`, http.StatusInternalServerError)
		return
	}
	expense.UserID = userId

	expenseRepo := repositories.NewExpenseRepository(config.GlobalConfig.DB)
	err = expenseRepo.Create(&expense)
	if err != nil {
		log.Println("Error adding expense: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, `{"error": "Error adding expense"}`, http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"message": "Expense added successfully",
		"link":    fmt.Sprintf("/dashboard/add-expense/%v", expense.ID),
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)

}

func Dashboard(w http.ResponseWriter, r *http.Request, limit int) {
	data := map[string]interface{}{}
	expenseRepo := repositories.NewExpenseRepository(config.GlobalConfig.DB)
	userId, err := utils.GetCookie(r, "usw")
	if err != nil {
		log.Println("Error getting user id: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, `{"error": "Error getting expenses"}`, http.StatusInternalServerError)
		return
	}
	expenses, err := expenseRepo.FindByUserID(userId, limit)
	if err != nil {
		log.Println("Error getting expenses: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, `{"error": "Error getting expenses"}`, http.StatusInternalServerError)
		return
	}
	data["Expenses"] = expenses
	templates.Render(w, "dashboard.html", data, true)
}
