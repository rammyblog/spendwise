package models

import "github.com/google/uuid"

type Category struct {
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primarykey" json:"id"`
	Name string    `json:"name"`
}
