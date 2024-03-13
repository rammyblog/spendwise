package models

import (
	"database/sql/driver"

	"github.com/google/uuid"
)

// UUID represents a UUID value
type UUID struct {
	uuid.UUID
}

// NewUUID creates a new UUID
func NewUUID() UUID {
	return UUID{UUID: uuid.New()}
}

// Scan scans a value into UUID
func (u *UUID) Scan(value interface{}) error {
	// Handle NULL values from the database
	if value == nil {
		u.UUID = uuid.UUID{}
		return nil
	}
	// Scan the value into UUID
	return u.UUID.Scan(value)
}

// Value returns the value of UUID
func (u UUID) Value() (driver.Value, error) {
	// Return UUID as a string
	return u.UUID.String(), nil
}
