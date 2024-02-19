package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/schema"
	"github.com/rammyblog/spendwise/config"
	"github.com/rammyblog/spendwise/models"
	"github.com/rammyblog/spendwise/repositories"
	"github.com/rammyblog/spendwise/utils"
)

func AddExpense(w http.ResponseWriter, r *http.Request) {
	var expense models.Expense

	err := r.ParseForm()
	if err != nil {
		log.Println("Error decoding form: ", err)
		utils.SetCookie(w, "errorSw", "Error adding expense", time.Now().Add(10*time.Second))
		http.Redirect(w, r, "/dashboard/add-expense", http.StatusSeeOther)
	}

	decoder := schema.NewDecoder()
	err = decoder.Decode(&expense, r.PostForm)
	if err != nil {
		log.Println("Error decoding form: ", err)
		utils.SetCookie(w, "errorSw", "Error adding expense", time.Now().Add(10*time.Second))
		http.Redirect(w, r, "/dashboard/add-expense", http.StatusSeeOther)
		return
	}

	userId, err := utils.GetCookie(r, "usw")

	if err != nil {
		log.Println("Error getting user id: ", err)
		utils.SetCookie(w, "errorSw", "Error adding expense", time.Now().Add(10*time.Second))
		http.Redirect(w, r, "/dashboard/add-expense", http.StatusSeeOther)
		return
	}
	expense.UserID = userId

	expenseRepo := repositories.NewExpenseRepository(config.GlobalConfig.DB)
	err = expenseRepo.Create(&expense)
	if err != nil {
		log.Println("Error adding expense: ", err)
		utils.SetCookie(w, "errorSw", "Error adding expense", time.Now().Add(10*time.Second))
		http.Redirect(w, r, "/dashboard/add-expense", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)

}
