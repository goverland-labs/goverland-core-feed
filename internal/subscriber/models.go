package subscriber

import "time"

type Subscriber struct {
	ID         string `gorm:"primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
	WebhookURL string
}
