package model

import "github.com/EdlanioJ/kbu-store/domain"

type Tag struct {
	Name  string `gorm:"column:tag;type:varchar"`
	Count int64  `gorm:"column:count;type:int"`
}

func (t *Tag) ToTagDomain() (res *domain.Tag) {
	res = new(domain.Tag)

	res.Count = t.Count
	res.Name = t.Name
	return
}
