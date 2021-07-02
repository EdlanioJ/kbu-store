package factory

import (
	"time"

	"github.com/EdlanioJ/kbu-store/data/usecase"
	gormRepo "github.com/EdlanioJ/kbu-store/infra/db/repository/gorm"
	"gorm.io/gorm"
)

func CategoryUsecase(db *gorm.DB, contextTimeout time.Duration) *usecase.CategoryUsecase {
	categoryRepo := gormRepo.NewCategoryRepository(db)

	return usecase.NewCategoryUsecase(categoryRepo, contextTimeout)
}
