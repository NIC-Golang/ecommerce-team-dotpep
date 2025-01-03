package models

type User struct {
	ID         string `json:"id,omitempty"`
	Email      string `json:"email,omitempty"`
	WebhookURL string `json:"webhook_url"`
	Phone      string `json:"phone,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
}
