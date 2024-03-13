package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/gorilla/schema"
	"github.com/rammyblog/spendwise/config"
	"github.com/rammyblog/spendwise/middleware"
	"github.com/rammyblog/spendwise/models"
	"github.com/rammyblog/spendwise/repositories"
	"github.com/rammyblog/spendwise/templates"
)

type CategoriesString struct {
	ID   string
	Name string
}

func AddExpense(w http.ResponseWriter, r *http.Request) {
	var expense models.Expense
	const layout = "2006-01-02"

	parsedDate, err := time.Parse(layout, r.FormValue("expense_date"))
	if err != nil {
		middleware.HandleError(w, err, "Error parsing date")
		return
	}
	expense.ExpenseDate = parsedDate

	err = r.ParseForm()

	if err != nil {
		middleware.HandleError(w, err, "Error adding expense")
	}

	r.PostForm.Set("expense_date", parsedDate.Format("2006-01-02T15:04:05Z07:00"))

	decoder := schema.NewDecoder()
	err = decoder.Decode(&expense, r.PostForm)
	if err != nil {
		fmt.Println("Error decoding form: ", err)
		middleware.HandleError(w, err, "Error adding expense")
		return
	}

	userId := r.Context().Value(middleware.UserIDKey).(string)

	expense.UserID = userId

	expenseRepo := repositories.NewExpenseRepository(config.GlobalConfig.DB)
	err = expenseRepo.Create(&expense)
	if err != nil {
		middleware.HandleError(w, err, "Error adding expense")
	}

	w.WriteHeader(http.StatusCreated)

	data := map[string]interface{}{}
	expensesPage := r.Header.Get("Expenses-Page")

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
	fmt.Println("userId: ", userId)
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

func ExpenseDetail(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	expenseRepo := repositories.NewExpenseRepository(config.GlobalConfig.DB)
	categtoryRepo := repositories.NewCategoryRepository(config.GlobalConfig.DB)
	expenseId := chi.URLParam(r, "id")
	expense, err := expenseRepo.FindByID(expenseId)
	if err != nil {
		templates.Render(w, "404.html", data, true)
		return
	}
	if expense.ID.String() == "" {
		templates.Render(w, "404.html", data, true)
		return
	}
	category, err := categtoryRepo.FindByID(expense.CategoryID)

	if err != nil {
		templates.Render(w, "404.html", data, true)
		return
	}
	expense.ExpenseDate = expense.ExpenseDate.Local()
	data["Expense"] = expense
	data["Category"] = category
	templates.Render(w, "expense-detail.html", data, true)
}

func EditExpenseForm(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	expenseRepo := repositories.NewExpenseRepository(config.GlobalConfig.DB)
	categtoryRepo := repositories.NewCategoryRepository(config.GlobalConfig.DB)
	expenseId := chi.URLParam(r, "id")
	expense, err := expenseRepo.FindByID(expenseId)
	if err != nil {
		templates.Render(w, "404.html", data, true)
		return
	}
	if expense.ID.String() == "" {
		templates.Render(w, "404.html", data, true)
		return
	}
	var categoriesString []CategoriesString
	categories, err := categtoryRepo.FindAll()

	// convert category id to string
	for _, category := range categories {
		categoriesString = append(categoriesString, CategoriesString{
			ID:   category.ID.String(),
			Name: category.Name,
		})
	}
	if err != nil {
		middleware.HandleError(w, err, "Error getting categories")
		return
	}

	data["Categories"] = categoriesString
	parsedDate := expense.ExpenseDate.Format("2006-01-02")
	data["Expense"] = expense
	data["ExpenseDate"] = parsedDate
	data["Link"] = "/dashboard/edit-expense/" + expenseId
	expensesDetailsPage := r.Header.Get("HX-Trigger")

	data["Target"] = "#expenses"
	if expensesDetailsPage == "detail-edit-btn" {
		data["Header"] = `{"Detail-Page": "true"}`
		data["Target"] = "#expense-details"
		data["State"] = "Update"
	}
	templates.Render(w, "add-expense.html", data, false)
}

func UpdateExpense(w http.ResponseWriter, r *http.Request) {
	var expense models.Expense
	const layout = "2006-01-02"
	parsedDate, err := time.Parse(layout, r.FormValue("expense_date"))
	if err != nil {
		middleware.HandleError(w, err, "Error parsing date")
		return
	}
	expense.ExpenseDate = parsedDate

	err = r.ParseForm()

	if err != nil {
		middleware.HandleError(w, err, "Error updating expense")
	}

	r.PostForm.Set("expense_date", parsedDate.Format("2006-01-02T15:04:05Z07:00"))

	decoder := schema.NewDecoder()
	err = decoder.Decode(&expense, r.PostForm)
	if err != nil {
		fmt.Println("Error decoding form: ", err)
		middleware.HandleError(w, err, "Error updating expense")
		return
	}

	expenseRepo := repositories.NewExpenseRepository(config.GlobalConfig.DB)
	expenseId := chi.URLParam(r, "id")
	userId := r.Context().Value(middleware.UserIDKey).(string)
	err = expenseRepo.Update(expenseId, &expense)
	if err != nil {
		middleware.HandleError(w, err, "Error updating expense")
		return
	}

	w.WriteHeader(http.StatusCreated)

	data := map[string]interface{}{}

	expenseDetailPage := r.Header.Get("Detail-Page")

	if expenseDetailPage == "true" {
		categoryRepo := repositories.NewCategoryRepository(config.GlobalConfig.DB)
		category, err := categoryRepo.FindByID(expense.CategoryID)
		if err != nil {
			middleware.HandleError(w, err, "Error updating expense")
			return
		}
		fmt.Println("category: ", expense)
		expense.ExpenseDate = expense.ExpenseDate.Local()
		expense.ID, _ = uuid.Parse(expenseId)
		data["Expense"] = expense
		data["Category"] = category
		templates.Render(w, "expense-details-partial.html", data, false)
		return
	}

	expenses, err := expenseRepo.FindByUserIDAndJoinCategory(userId, 10)
	if err != nil {
		middleware.HandleError(w, err, "Error getting expenses")
		return
	}

	data["Expenses"] = expenses
	templates.Render(w, "expense-table.html", data, false)
}

func DeleteExpense(w http.ResponseWriter, r *http.Request) {
	expenseRepo := repositories.NewExpenseRepository(config.GlobalConfig.DB)
	expenseId := chi.URLParam(r, "id")
	userId := r.Context().Value(middleware.UserIDKey).(string)
	// confirm user owns the expense
	expense, err := expenseRepo.FindByID(expenseId)
	if err != nil {
		middleware.HandleError(w, err, "Error deleting expense")
		return
	}
	if expense.UserID != userId {
		middleware.HandleError(w, err, "Error deleting expense")
		return
	}
	err = expenseRepo.Delete(expenseId)
	if err != nil {
		middleware.HandleError(w, err, "Error deleting expense")
		return
	}

	data := map[string]interface{}{}

	expenses, err := expenseRepo.FindByUserIDAndJoinCategory(userId, 10)
	if err != nil {
		middleware.HandleError(w, err, "Error getting expenses")
		return
	}

	// check if the request is from the expenses page
	expensesDetailsPage := r.Header.Get("HX-Trigger")

	if expensesDetailsPage == "detail-delete-btn" {
		url := "http://" + r.Host + "/dashboard/expenses"
		w.Header().Set("HX-Redirect", url)
		return
	}

	data["Expenses"] = expenses
	templates.Render(w, "expense-table.html", data, false)

}
