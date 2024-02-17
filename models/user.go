package models

type User struct {
	Model
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	Picture       string `json:"picture"`
	Provider      string `json:"provider"`
	EmailVerified bool   `json:"email_verified"`
	ProviderID    string `json:"provider_id"`
}
