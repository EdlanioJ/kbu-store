package domain

type Position struct {
	Lat float64 `json:"lat" valid:"latitude,optional"`
	Lng float64 `json:"lng" valid:"longitude,optional"`
}
