package models

import "github.com/google/uuid"

type Expense struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primarykey" json:"id"`
	Amount      float64   `json:"amount"`
	UserID      string    `json:"user_id"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	ExpenseDate string    `json:"expense_date"`
	CreatedAt   string    `json:"created_at"`
}
