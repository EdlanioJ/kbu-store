package model

import "github.com/EdlanioJ/kbu-store/app/domain"

type Category struct {
	Base
	Name   string   `gorm:"column:name;type:varchar;not null"`
	Status string   `gorm:"type:varchar(20)"`
	Store  []*Store `gorm:"ForeignKey:CategoryID"`
}

func (Category) TableName() string {
	return "categories"
}

func (s *Category) FromCategoryDomain(d *domain.Category) *Category {
	s.ID = d.ID
	s.Name = d.Name
	s.Status = d.Status
	s.CreatedAt = d.CreatedAt
	s.UpdatedAt = d.UpdatedAt

	return s
}

func (s *Category) ToCategoryDomain() (res *domain.Category) {
	res = new(domain.Category)
	res.ID = s.ID
	res.Name = s.Name
	res.Status = s.Status
	res.CreatedAt = s.CreatedAt
	res.UpdatedAt = s.UpdatedAt

	return
}
