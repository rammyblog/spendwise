package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/rammyblog/spendwise/config"
	"github.com/rammyblog/spendwise/middleware"
	"github.com/rammyblog/spendwise/models"
	"github.com/rammyblog/spendwise/repositories"
	"github.com/rammyblog/spendwise/templates"
)

func AddExpense(w http.ResponseWriter, r *http.Request) {
	var expense models.Expense

	err := r.ParseForm()
	if err != nil {
		middleware.HandleError(w, err, "Error adding expense")
	}

	decoder := schema.NewDecoder()
	err = decoder.Decode(&expense, r.PostForm)
	if err != nil {
		middleware.HandleError(w, err, "Error adding expense")

	}

	userId := r.Context().Value(middleware.UserIDKey).(string)

	if err != nil {
		middleware.HandleError(w, err, "Error adding expense")

	}
	expense.UserID = userId

	expenseRepo := repositories.NewExpenseRepository(config.GlobalConfig.DB)
	err = expenseRepo.Create(&expense)
	if err != nil {
		middleware.HandleError(w, err, "Error adding expense")

	}

	w.WriteHeader(http.StatusCreated)

	data := map[string]interface{}{}
	expensesPage := r.Header.Get("Expenses-Page")
	fmt.Println(expensesPage, "expensesPage")

	if expensesPage == "true" {
		expenses, err := expenseRepo.FindByUserIDAndJoinCategory(userId, 10)
		if err != nil {
			middleware.HandleError(w, err, "Error getting expenses")
		}
		data["Expenses"] = expenses
		templates.Render(w, "expense-table.html", data, false)
		return
	}

	expenses, err := expenseRepo.FindByUserID(userId, 5)
	if err != nil {
		middleware.HandleError(w, err, "Error adding expense")
	}

	maxAmountForCategory, expenseForAMonth, totalExpenses, totalAmountPerCategory := expenseRepo.GetExpenseSummary(userId)
	data["MaxAmountForCategory"] = maxAmountForCategory
	data["ExpenseForAMonth"] = expenseForAMonth
	data["TotalExpenses"] = totalExpenses
	data["Expenses"] = expenses
	data["MaxCategory"] = maxAmountForCategory.CategoryName
	data["MaxAmount"] = maxAmountForCategory.Amount
	data["ChartData"] = totalAmountPerCategory

	templates.Render(w, "expense-stats-grid.html", data, false)
}

func ExpenseGraph(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middleware.UserIDKey).(string)

	expenseRepo := repositories.NewExpenseRepository(config.GlobalConfig.DB)
	_, _, _, totalAmountPerCategory := expenseRepo.GetExpenseSummary(userId)
	data := map[string]interface{}{
		"ChartData": totalAmountPerCategory,
	}
	templates.Render(w, "pie.html", data, false)
}

func Dashboard(w http.ResponseWriter, r *http.Request, limit int) {
	data := map[string]interface{}{}
	expenseRepo := repositories.NewExpenseRepository(config.GlobalConfig.DB)
	userId := r.Context().Value(middleware.UserIDKey).(string)

	expenses, err := expenseRepo.FindByUserID(userId, limit)
	if err != nil {
		middleware.HandleError(w, err, "Error getting expenses")
	}
	maxAmountForCategory, expenseForAMonth, totalExpenses, totalAmountPerCategory := expenseRepo.GetExpenseSummary(userId)
	data["MaxAmountForCategory"] = maxAmountForCategory
	data["ExpenseForAMonth"] = expenseForAMonth
	data["TotalExpenses"] = totalExpenses
	data["Expenses"] = expenses
	data["MaxCategory"] = maxAmountForCategory.CategoryName
	data["MaxAmount"] = maxAmountForCategory.Amount
	data["ChartData"] = totalAmountPerCategory

	templates.Render(w, "dashboard.html", data, true)
}

func ExpenseList(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	expenseRepo := repositories.NewExpenseRepository(config.GlobalConfig.DB)
	userId := r.Context().Value(middleware.UserIDKey).(string)

	expenses, err := expenseRepo.FindByUserIDAndJoinCategory(userId, 10)
	if err != nil {
		middleware.HandleError(w, err, "Error getting expenses")
	}
	data["Expenses"] = expenses
	templates.Render(w, "expenses.html", data, true)
}
