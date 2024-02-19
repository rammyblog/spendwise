package repositories

import (
	"github.com/google/uuid"
	"github.com/rammyblog/spendwise/models"
	"gorm.io/gorm"
)

type ExpenseRepository struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (repo *ExpenseRepository) Create(expense *models.Expense) error {
	return repo.db.Create(expense).Error
}

func (repo *ExpenseRepository) FindByID(id uuid.UUID) (*models.Expense, error) {
	var expense models.Expense
	err := repo.db.First(&expense, id).Error
	return &expense, err
}

func (repo *ExpenseRepository) FindAll() ([]models.Expense, error) {
	var expenses []models.Expense
	err := repo.db.Find(&expenses).Error
	return expenses, err
}

func (repo *ExpenseRepository) Update(id uuid.UUID, expense *models.Expense) error {
	var existingExpense models.Expense
	err := repo.db.First(&existingExpense, id).Error
	if err != nil {
		return err
	}
	return repo.db.Model(&existingExpense).Updates(expense).Error
}

func (repo *ExpenseRepository) Delete(id uuid.UUID) error {
	var expense models.Expense
	err := repo.db.First(&expense, id).Error
	if err != nil {
		return err
	}
	return repo.db.Delete(&expense).Error
}

func (repo *ExpenseRepository) FindByUserID(userID uuid.UUID) ([]models.Expense, error) {
	var expenses []models.Expense
	err := repo.db.Where("user_id = ?", userID).Find(&expenses).Error
	return expenses, err
}

func (repo *ExpenseRepository) FindByCategory(categoryId uuid.UUID) ([]models.Expense, error) {
	var expenses []models.Expense
	err := repo.db.Where("category_id = ?", categoryId).Find(&expenses).Error
	return expenses, err
}
