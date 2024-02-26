package repositories

import (
	"time"

	"github.com/rammyblog/spendwise/models"
	"gorm.io/gorm"
)

type ExpenseWithCategory struct {
	models.Expense
	CategoryName string `json:"category_name"`
}

type MaxAmountForCategory struct {
	CategoryName string
	Amount       float64
}

type TotalAmountForCategory struct {
	CategoryName string  `json:"category_name"`
	TotalAmount  float64 `json:"total_amount"`
	CategoryID   string  `json:"category_id"`
}

type ExpenseRepository struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (repo *ExpenseRepository) Create(expense *models.Expense) error {
	return repo.db.Create(expense).Error
}

func (repo *ExpenseRepository) FindByID(id string) (*models.Expense, error) {
	var expense models.Expense
	err := repo.db.First(&expense, id).Error
	return &expense, err
}

func (repo *ExpenseRepository) FindAll() ([]models.Expense, error) {
	var expenses []models.Expense
	err := repo.db.Find(&expenses).Error
	return expenses, err
}

func (repo *ExpenseRepository) Update(id string, expense *models.Expense) error {
	var existingExpense models.Expense
	err := repo.db.First(&existingExpense, id).Error
	if err != nil {
		return err
	}
	return repo.db.Model(&existingExpense).Updates(expense).Error
}

func (repo *ExpenseRepository) Delete(id string) error {
	var expense models.Expense
	err := repo.db.First(&expense, id).Error
	if err != nil {
		return err
	}
	return repo.db.Delete(&expense).Error
}

func (repo *ExpenseRepository) FindByUserID(userID string, limit int) ([]models.Expense, error) {
	var expenses []models.Expense
	err := repo.db.Where("user_id = ?", userID).Limit(limit).Find(&expenses).Error
	return expenses, err
}

func (repo *ExpenseRepository) FindByUserIDAndJoinCategory(userID string, limit int) ([]ExpenseWithCategory, error) {
	var expensesWithCategory []ExpenseWithCategory
	// err := repo.db.Joins("Category").Where("user_id = ?", userID).Limit(limit).Find(&expensesWithCategory).Error
	err := repo.db.Table("expenses").
		Select("expenses.name, expenses.amount, expenses.id, expenses.expense_date, categories.name as category_name").
		Joins("left join categories on categories.id = expenses.category_id").
		Where("user_id = ?", userID).
		Limit(limit).
		Order("expenses.created_at desc").
		Scan(&expensesWithCategory).Error

	return expensesWithCategory, err
}

func (repo *ExpenseRepository) FindByCategory(categoryId string) ([]models.Expense, error) {
	var expenses []models.Expense
	err := repo.db.Where("category_id = ?", categoryId).Find(&expenses).Error
	return expenses, err
}

func (repo *ExpenseRepository) GetExpenseSummary(userID string) (MaxAmountForCategory, float64, float64, []TotalAmountForCategory) {

	var expenseForAMonth float64
	var totalExpenses float64

	var maxAmountForCategory MaxAmountForCategory
	var totalAmountPerCategory []TotalAmountForCategory

	sqlForMaxAmount := `
	SELECT e.id AS expense_id, e.amount, e.category_id, c.name AS category_name, e.user_id
	FROM expenses e
	JOIN (
		SELECT category_id, MAX(amount) AS max_amount
		FROM expenses
		WHERE user_id = ?
		GROUP BY category_id
		ORDER BY max_amount DESC
		LIMIT 1
	) t ON e.category_id = t.category_id AND e.amount = t.max_amount
	JOIN categories c ON e.category_id = c.id
`

	sqlForCategoryAndExpenses := `	SELECT c.id AS category_id, c.name AS category_name, SUM(e.amount) AS total_amount
	FROM expenses e
	JOIN categories c ON e.category_id = c.id
	WHERE e.user_id = ?
	GROUP BY c.id, c.name
	ORDER BY total_amount DESC
`

	repo.db.Raw(sqlForMaxAmount, userID).Scan(&maxAmountForCategory)
	repo.db.Raw(sqlForCategoryAndExpenses, userID).Scan(&totalAmountPerCategory)

	// Get the current time
	now := time.Now()

	// Get the first day of the current month
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// Get the first day of the next month
	firstOfNextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
	repo.db.Table("expenses").Select("sum(amount)").Where("expense_date >= ? AND expense_date < ? AND user_id = ?", firstOfMonth, firstOfNextMonth, userID).Row().Scan(&expenseForAMonth)
	repo.db.Table("expenses").Select("sum(amount)").Where("user_id = ?", userID).Row().Scan(&totalExpenses)

	return maxAmountForCategory, expenseForAMonth, totalExpenses, totalAmountPerCategory
}
