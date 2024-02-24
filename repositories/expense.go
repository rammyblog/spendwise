package repositories

import (
	"time"

	"github.com/rammyblog/spendwise/models"
	"gorm.io/gorm"
)

type MaxAmountForCategory struct {
	Category string
	Amount   float64
	UserID   string
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

func (repo *ExpenseRepository) FindByCategory(categoryId string) ([]models.Expense, error) {
	var expenses []models.Expense
	err := repo.db.Where("category_id = ?", categoryId).Find(&expenses).Error
	return expenses, err
}

func (repo *ExpenseRepository) GetExpenseSummary(userID string) (MaxAmountForCategory, float64, float64) {

	var expenseForAMonth float64
	var totalExpenses float64

	var result MaxAmountForCategory
	subQuery := repo.db.Table("expenses").
		Select("MAX(amount) as amount").
		Where("user_id = ?", userID)

	// Define the main query
	repo.db.Table("expenses").
		Select("amount, category").
		Where("amount = (?) AND user_id = ?", subQuery, userID).
		First(&result)
	// Get the current time
	now := time.Now()

	// Get the first day of the current month
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// Get the first day of the next month
	firstOfNextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())

	repo.db.Table("expenses").Select("sum(amount)").Where("expense_date >= ? AND expense_date < ? AND user_id = ?", firstOfMonth, firstOfNextMonth, userID).Row().Scan(&expenseForAMonth)
	repo.db.Table("expenses").Select("sum(amount)").Where("user_id = ?", userID).Row().Scan(&totalExpenses)

	return result, expenseForAMonth, totalExpenses
}
