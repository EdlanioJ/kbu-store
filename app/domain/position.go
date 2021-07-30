package domain

type Position struct {
	Lat float64 `json:"lat" gorm:"type:decimal(11,8)"`
	Lng float64 `json:"lng" gorm:"type:decimal(11,8)"`
}
