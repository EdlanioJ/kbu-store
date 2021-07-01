package factory

import (
	"time"

	gormRepo "github.com/EdlanioJ/kbu-store/category/repository/gorm"
	"github.com/EdlanioJ/kbu-store/data/usecase"
	"gorm.io/gorm"
)

func CategoryUsecase(db *gorm.DB, contextTimeout time.Duration) *usecase.CategoryUsecase {
	categoryRepo := gormRepo.NewCategoryRepository(db)

	return usecase.NewCategoryUsecase(categoryRepo, contextTimeout)
}
