package model

import "github.com/EdlanioJ/kbu-store/domain"

type Category struct {
	Base
	Name   string   `gorm:"column:name;type:varchar;not null"`
	Status string   `gorm:"type:varchar(20)"`
	Store  []*Store `gorm:"ForeignKey:CategoryID"`
}

func (Category) TableName() string {
	return "categories"
}

func (s *Category) Parser(d *domain.Category) *Category {
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
