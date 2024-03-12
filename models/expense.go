package models

import (
	"time"

	"github.com/google/uuid"
)

type Expense struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primarykey" schema:"id"`
	Amount      float64   `schema:"amount"`
	UserID      string    `schema:"user_id"`
	Name        string    `schema:"name"`
	Description string    `schema:"description"`
	CategoryID  string    `schema:"category_id"`
	ExpenseDate time.Time `schema:"expense_date"`
	CreatedAt   time.Time `schema:"created_at"`
}
