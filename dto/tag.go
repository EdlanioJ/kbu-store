package dto

import "github.com/EdlanioJ/kbu-store/domain"

type TagDBModel struct {
	Name  string `gorm:"column:tag;type:varchar"`
	Count int64  `gorm:"column:count;type:int"`
}

func (t *TagDBModel) ParserToTagDomain() (res *domain.Tag) {
	res = new(domain.Tag)

	res.Count = t.Count
	res.Name = t.Name
	return
}
