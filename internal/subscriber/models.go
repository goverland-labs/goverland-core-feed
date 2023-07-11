package subscriber

import (
	"time"

	"gorm.io/gorm"
)

type Subscriber struct {
	ID         string `gorm:"primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	WebhookURL string
}
