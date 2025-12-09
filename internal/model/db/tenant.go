package model

type Tenant struct {
	ID         uint   `gorm:"primaryKey"`
	AccountID  string `gorm:"size:24"`
	WebhookURL string `gorm:"size:256"`
	APIKey     string `gorm:"size:12"`
	Name       string `gorm:"size:64"`
}
