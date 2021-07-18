package handler

import "github.com/EdlanioJ/kbu-store/domain"

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type CreateStoreRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	CategoryID  string   `json:"category_id"`
	UserID      string   `json:"user_id"`
	Tags        []string `json:"tags"`
	Lat         float64  `json:"latitude"`
	Lng         float64  `json:"longitude"`
}

type UpdateStoreRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	CategoryID  string   `json:"category_id"`
	Tags        []string `json:"tags"`
	Lat         float64  `json:"latitude"`
	Lng         float64  `json:"longitude"`
}

func (r *UpdateStoreRequest) ToDomainStore() (res *domain.Store) {
	res = new(domain.Store)
	category := new(domain.Category)
	category.ID = r.CategoryID

	res.Name = r.Name
	res.Description = r.Description
	res.Tags = r.Tags
	res.Position.Lat = r.Lat
	res.Position.Lng = r.Lng
	res.Category = category

	return
}
