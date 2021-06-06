package factory

import (
	"time"

	"github.com/EdlanioJ/kbu-store/tag/handler"
	tagRepo "github.com/EdlanioJ/kbu-store/tag/repository/gorm"
	"github.com/EdlanioJ/kbu-store/tag/usecase"
	"gorm.io/gorm"
)

func FetchTagsHandler(db *gorm.DB, contextTimeout time.Duration) *handler.FetchTags {
	fetchTagsRepo := tagRepo.NewGormFetchTags(db)
	fetchTagsUsecase := usecase.NewFetchTags(fetchTagsRepo, contextTimeout)

	return handler.NewFetchTags(fetchTagsUsecase)
}

func FetchTagsByCategoryHandler(db *gorm.DB, contextTimeout time.Duration) *handler.FetchTagsByCategory {
	fetchTagsRepo := tagRepo.NewGormFetchTagsByCategory(db)
	fetchTagsUsecase := usecase.NewFetchTagsByCategory(fetchTagsRepo, contextTimeout)

	return handler.NewFetchTagsByCategory(fetchTagsUsecase)
}
