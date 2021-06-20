package domain

import "context"

type Tag struct {
	Name  string `json:"tag"`
	Count int64  `json:"count,omitempty"`
}

type (
	// TagRepository
	TagRepository interface {
		GetAll(ctx context.Context, sort string, page, limit int) ([]*Tag, int64, error)
		GetAllByCategory(ctx context.Context, category, sort string, page, limit int) ([]*Tag, int64, error)
	}

	// TagUsecase
	TagUsecase interface {
		GetAll(ctx context.Context, sort string, page, limit int) ([]*Tag, int64, error)
		GetAllByCategory(ctx context.Context, categoryID, sort string, page, limit int) ([]*Tag, int64, error)
	}
)
