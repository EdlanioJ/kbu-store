package model

import (
	"github.com/EdlanioJ/kbu-store/app/domain"
	"github.com/lib/pq"
)

type Store struct {
	Base
	Name        string         `gorm:"column:name;type:varchar;not null"`
	Description string         `gorm:"type:varchar(255)"`
	Status      string         `gorm:"type:varchar(20)"`
	UserID      string         `gorm:"column:user_id;type:uuid"`
	AccountID   string         `gorm:"column:account_id;type:uuid"`
	CategoryID  string         `gorm:"column:category_id;type:uuid"`
	Tags        pq.StringArray `gorm:"type:text[]"`
	Lat         float64        `gorm:"type:decimal(10,8)"`
	Lng         float64        `gorm:"type:decimal(11,8)"`
}

func (Store) TableName() string {
	return "stores"
}

func (s *Store) FromStoreDomain(d *domain.Store) {
	s.ID = d.ID
	s.Name = d.Name
	s.Description = d.Description
	s.UserID = d.UserID
	s.Status = d.Status
	s.AccountID = d.AccountID
	s.CategoryID = d.CategoryID
	s.Tags = pq.StringArray(d.Tags)
	s.Lat = d.Position.Lat
	s.Lng = d.Position.Lng
	s.CreatedAt = d.CreatedAt
	s.UpdatedAt = d.UpdatedAt
}

func (s *Store) ToStoreDomain() (res *domain.Store) {
	res = new(domain.Store)
	account := new(domain.Account)

	account.ID = s.AccountID

	res.ID = s.ID
	res.Name = s.Name
	res.Description = s.Description
	res.UserID = s.UserID
	res.Status = s.Status
	res.Tags = s.Tags
	res.Position.Lat = s.Lat
	res.Position.Lng = s.Lng
	res.AccountID = s.AccountID
	res.CategoryID = s.CategoryID
	res.CreatedAt = s.CreatedAt
	res.UpdatedAt = s.UpdatedAt

	return
}