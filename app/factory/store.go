package factory

import (
	"time"

	"github.com/EdlanioJ/kbu-store/data/usecase"
	gormRepo "github.com/EdlanioJ/kbu-store/infra/db/repository/gorm"
	gormStoreRepo "github.com/EdlanioJ/kbu-store/store/repository/gorm"
	"gorm.io/gorm"
)

func StoreUsecase(db *gorm.DB, contextTimeout time.Duration) *usecase.StoreUsecase {
	storeRepo := gormStoreRepo.NewStoreRepository(db)
	accountRepo := gormRepo.NewAccountRepository(db)
	categoryRepo := gormRepo.NewCategoryRepository(db)

	return usecase.NewStoreUsecase(storeRepo, accountRepo, categoryRepo, contextTimeout)
}

func TagUsecase(db *gorm.DB, contextTimeout time.Duration) *usecase.TagUsecase {
	tagRepo := gormStoreRepo.NewTagsRepository(db)

	return usecase.NewtagUsecase(tagRepo, contextTimeout)
}
