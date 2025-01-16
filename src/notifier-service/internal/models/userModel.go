package models

type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	WebhookURL string `json:"webhook_url"`
	Phone      string `json:"phone"`
	CreatedAt  string `json:"created_at"`
	UserType   string `json:"user_type"`
	Timezone   string `json:"timezone"`
}
