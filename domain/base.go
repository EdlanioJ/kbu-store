package domain

import (
	"time"

	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Base struct {
	ID        string    `json:"id" valid:"uuidv4"`
	CreatedAt time.Time `json:"created_at,omitempty" valid:"-"`
	UpdatedAt time.Time `json:"-" valid:"-"`
}
