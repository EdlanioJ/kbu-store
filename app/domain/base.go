package domain

import (
	"time"
)

type Base struct {
	ID        string    `json:"id" gorm:"column:id;type:uuid;primary key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}
