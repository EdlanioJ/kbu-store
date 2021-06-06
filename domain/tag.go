package domain

import "context"

type Tag struct {
	Name  string `json:"tag"`
	Count int64  `json:"count,omitempty"`
}

type (
	FetchTagsRepository interface {
		Exec(ctx context.Context, sort string, page, limit int) ([]*Tag, int64, error)
	}
	FetchTagsByCategoryRepository interface {
		Exec(ctx context.Context, category, sort string, page, limit int) ([]*Tag, int64, error)
	}
)

type (
	FetchTagsUsecase interface {
		Exec(ctx context.Context, sort string, page, limit int) ([]*Tag, int64, error)
	}
	FetchTagsByCategoryUsecase interface {
		Exec(ctx context.Context, category, sort string, page, limit int) ([]*Tag, int64, error)
	}
)
