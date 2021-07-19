package factory

import (
	"time"

	gormRepo "github.com/EdlanioJ/kbu-store/infra/db/repository/gorm"
	"github.com/EdlanioJ/kbu-store/usecase"
	"gorm.io/gorm"
)

func StoreUsecase(db *gorm.DB, contextTimeout time.Duration) *usecase.StoreUsecase {
	storeRepo := gormRepo.NewStoreRepository(db)
	accountRepo := gormRepo.NewAccountRepository(db)
	categoryRepo := gormRepo.NewCategoryRepository(db)

	return usecase.NewStoreUsecase(storeRepo, accountRepo, categoryRepo, contextTimeout)
}
