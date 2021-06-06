package dto

import "github.com/EdlanioJ/kbu-store/domain"

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type CreateStoreRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	CategoryID  string   `json:"category_id"`
	ExternalID  string   `json:"external_id"`
	Tags        []string `json:"tags"`
	Lat         float64  `json:"latitude"`
	Lng         float64  `json:"longitude"`
}

type TagsQueryParams struct {
	Page  int      `query:"page"`
	Limit int      `query:"limit"`
	Sort  string   `query:"sort"`
	Tags  []string `query:"tags"`
}

type UpdateStoreRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	CategoryID  string   `json:"category_id"`
	Tags        []string `json:"tags"`
	Lat         float64  `json:"latitude"`
	Lng         float64  `json:"longitude"`
}

func (r *UpdateStoreRequest) Parser() (res *domain.Store) {
	res = new(domain.Store)
	Category := new(domain.Category)
	Category.ID = r.CategoryID

	res.Name = r.Name
	res.Description = r.Description
	res.Tags = r.Tags
	res.Position.Lat = r.Lat
	res.Position.Lng = r.Lng
	res.Category = Category

	return
}
