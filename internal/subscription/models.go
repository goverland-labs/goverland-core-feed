package subscription

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subscription struct {
	ID           uuid.UUID `gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	SubscriberID uuid.UUID
	DaoID        uuid.UUID
}
