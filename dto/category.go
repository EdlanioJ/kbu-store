package dto

import "github.com/EdlanioJ/kbu-store/domain"

type CategoryDBModel struct {
	Base
	Name   string          `gorm:"column:name;type:varchar;not null"`
	Status string          `gorm:"type:varchar(20)"`
	Store  []*StoreDBModel `gorm:"ForeignKey:CategoryID"`
}

func (CategoryDBModel) TableName() string {
	return "categories"
}

func (s *CategoryDBModel) ParserToDBModel(d *domain.Category) *CategoryDBModel {
	s.ID = d.ID
	s.Name = d.Name
	s.Status = d.Status
	s.CreatedAt = d.CreatedAt
	s.UpdatedAt = d.UpdatedAt

	return s
}

func (s *CategoryDBModel) ParserToCategoryDomain() (res *domain.Category) {
	res = new(domain.Category)
	res.ID = s.ID
	res.Name = s.Name
	res.Status = s.Status
	res.CreatedAt = s.CreatedAt
	res.UpdatedAt = s.UpdatedAt

	return
}
