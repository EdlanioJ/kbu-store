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
	CategoryStatusInactive string = "inactive"
)

// Category struct
type Category struct {
	Base   `valid:"required"`
	Name   string `json:"name" valid:"notnull"`
	Status string `json:"status" valid:"notnull,status"`
}

// Repositories
type (
	FetchCategoryRepository interface {
		Exec(ctx context.Context, sort string, page, limit int) ([]*Category, int64, error)
	}
	FetchCategoryByStatusRepository interface {
		Exec(ctx context.Context, status, sort string, page, limit int) ([]*Category, int64, error)
	}
	GetCategoryByIDRepository interface {
		Exec(ctx context.Context, id string) (*Category, error)
	}
	GetCategoryByStautsRepository interface {
		Exec(ctx context.Context, id, status string) (*Category, error)
	}
	CreateCategoryRepository interface {
		Add(ctx context.Context, Category *Category) error
	}
	UpdateCategoryRepository interface {
		Exec(ctx context.Context, Category *Category) error
	}
)

// Usecases
type (
	FetchCategoryUsecase interface {
		Exec(ctx context.Context, sort string, page, limit int) ([]*Category, int64, error)
	}
	FetchCategoryByStatusUsecase interface {
		Exec(ctx context.Context, status, sort string, page, limit int) ([]*Category, int64, error)
	}
	GetCategoryByIDUsecase interface {
		Exec(ctx context.Context, id string) (*Category, error)
	}
	GetCategoryByStautsUsecase interface {
		Exec(ctx context.Context, id, status string) (*Category, error)
	}
	CreateCategoryUsecase interface {
		Add(ctx context.Context, name string) error
	}
	UpdateCategoryUsecase interface {
		Exec(ctx context.Context, Category *Category) error
	}
)

// Category Entity validator
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
