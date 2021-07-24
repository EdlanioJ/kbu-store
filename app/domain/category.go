package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/asaskevich/govalidator"
)

const (
	CategoryStatusPending  string = "pending"
	CategoryStatusActive   string = "active"
	CategoryStatusInactive string = "disable"
)

// Category struct
type Category struct {
	Base   `valid:"required"`
	Name   string `json:"name" valid:"notnull"`
	Status string `json:"status" valid:"notnull,status"`
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

// Category validator
func (s *Category) isValid() (err error) {
	govalidator.TagMap["status"] = govalidator.Validator(func(str string) bool {
		return govalidator.IsIn(str, CategoryStatusPending, CategoryStatusActive, CategoryStatusInactive)
	})

	_, err = govalidator.ValidateStruct(s)

	return
}

// NewCategory creates an *Category struct
func NewCategory(id, name, status string) (category *Category, err error) {
	category = &Category{
		Name: name,
	}

	category.ID = id
	category.Status = status
	category.CreatedAt = time.Now()

	err = category.isValid()

	if err != nil {
		return nil, err
	}

	return
}

func (c *Category) ParseJson(data []byte) (err error) {
	err = json.Unmarshal(data, c)
	if err != nil {
		return
	}

	err = c.isValid()
	return
}

func (c *Category) Activate() (err error) {
	c.Status = CategoryStatusActive
	err = c.isValid()
	return
}

func (c *Category) Disable() (err error) {
	c.Status = CategoryStatusInactive
	err = c.isValid()
	return
}
