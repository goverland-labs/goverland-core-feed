package subscription

import "gorm.io/gorm"

type Subscription struct {
	gorm.Model
	SubscriberID string
	DaoID        string
}
