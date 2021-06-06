package domain

type Position struct {
	Lat float64 `json:"latitude" valid:"latitude,optional"`
	Lng float64 `json:"longitude" valid:"longitude,optional"`
}
