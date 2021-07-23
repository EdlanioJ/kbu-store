package domain

import (
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
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
func NewCategory(name string) (category *Category, err error) {
	category = &Category{
		Name: name,
	}

	category.ID = uuid.NewV4().String()
	category.CreatedAt = time.Now()
	category.Status = CategoryStatusPending

	err = category.isValid()

	if err != nil {
		return nil, err
	}

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
