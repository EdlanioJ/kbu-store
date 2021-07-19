package factory

import (
	"time"

	gormRepo "github.com/EdlanioJ/kbu-store/infra/db/repository/gorm"
	"github.com/EdlanioJ/kbu-store/usecase"
	"gorm.io/gorm"
)

func CategoryUsecase(db *gorm.DB, contextTimeout time.Duration) *usecase.CategoryUsecase {
	categoryRepo := gormRepo.NewCategoryRepository(db)

	return usecase.NewCategoryUsecase(categoryRepo, contextTimeout)
}
