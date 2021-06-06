package dto

import (
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
)

type CategoryDBModel struct {
	ID        string `gorm:"column:id;type:uuid;primary key"`
	Name      string `gorm:"column:name;type:varchar;not null"`
	Status    string `gorm:"type:varchar(20)"`
	CreatedAt time.Time
	UpdatedAt time.Time
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
