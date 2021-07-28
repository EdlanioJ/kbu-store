package sample

import (
	"time"

	"github.com/EdlanioJ/kbu-store/app/domain"
	uuid "github.com/satori/go.uuid"
)

func NewCategory() *domain.Category {
	c := new(domain.Category)
	c.ID = uuid.NewV4().String()
	c.Name = "Category 001"
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	c.Status = domain.CategoryStatusPending
	return c
}
