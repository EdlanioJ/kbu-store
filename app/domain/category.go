package domain

import (
	"context"
	"encoding/json"
	"time"
)

const (
	CategoryStatusPending string = "pending"
	CategoryStatusActive  string = "active"
	CategoryStatusDisable string = "disable"
)

// Category struct
type Category struct {
	Base
	Name   string `json:"name" gorm:"column:name;type:varchar;not null"`
	Status string `json:"status" gorm:"type:varchar(20)"`
}

type (
	// CategoryRepository
	CategoryRepository interface {
		Store(ctx context.Context, Category *Category) error
		FindByID(ctx context.Context, id string) (*Category, error)
		Update(ctx context.Context, Category *Category) error
	}
	// CategoryUsecase
	CategoryUsecase interface {
		Create(ctx context.Context, Category *Category) error
		Update(ctx context.Context, Category *Category) error
	}
)

// NewCategory creates an *Category struct
func NewCategory() (category *Category) {
	category = new(Category)
	category.Status = CategoryStatusPending
	category.CreatedAt = time.Now()

	return
}

// Parses a JSON data and store in the Category Entity
func (c *Category) ParseJson(data []byte) (err error) {
	err = json.Unmarshal(data, c)
	if err != nil {
		return
	}
	return
}
