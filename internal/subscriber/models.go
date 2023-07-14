package subscriber

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subscriber struct {
	ID         uuid.UUID `gorm:"primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	WebhookURL string
}
