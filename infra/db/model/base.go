package model

import "time"

type Base struct {
	ID        string `gorm:"column:id;type:uuid;primary key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
