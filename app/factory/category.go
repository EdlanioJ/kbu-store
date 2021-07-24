package factory

import (
	"time"

	gormRepo "github.com/EdlanioJ/kbu-store/app/infrastructure/repository/gorm"
	"github.com/EdlanioJ/kbu-store/app/usecases"
	"gorm.io/gorm"
)

func CategoryUsecase(db *gorm.DB, contextTimeout time.Duration) *usecases.CategoryUsecase {
	categoryRepo := gormRepo.NewCategoryRepository(db)

	return usecases.NewCategoryUsecase(categoryRepo, contextTimeout)
}
