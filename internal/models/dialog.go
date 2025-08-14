package models

import (
	"time"
)

type Message struct {
	ID        int64 `gorm:"primaryKey"`
	CreatedAt time.Time
	UserID    int64 `gorm:"index"`
	Role      string
	Content   string
}
