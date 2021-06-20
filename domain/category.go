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
		Create(ctx context.Context, Category *Category) error
		GetById(ctx context.Context, id string) (*Category, error)
		GetByIdAndStatus(ctx context.Context, id, status string) (*Category, error)
		GetAll(ctx context.Context, sort string, page, limit int) ([]*Category, int64, error)
		GetAllByStatus(ctx context.Context, status, sort string, page, limit int) ([]*Category, int64, error)
		Update(ctx context.Context, Category *Category) error
	}
	// CategoryUsecase
	CategoryUsecase interface {
		Create(ctx context.Context, name string) error
		GetById(ctx context.Context, id string) (*Category, error)
		GetByIdAndStatus(ctx context.Context, id, status string) (*Category, error)
		GetAll(ctx context.Context, sort string, page, limit int) ([]*Category, int64, error)
		GetAllByStatus(ctx context.Context, status, sort string, page, limit int) ([]*Category, int64, error)
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
