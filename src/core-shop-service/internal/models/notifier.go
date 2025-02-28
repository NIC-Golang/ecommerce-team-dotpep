package models

type Notifier struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Description string `json:"description"`
}
