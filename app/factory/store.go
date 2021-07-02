package factory

import (
	"time"

	gormCategoryRepo "github.com/EdlanioJ/kbu-store/category/repository/gorm"
	"github.com/EdlanioJ/kbu-store/data/usecase"
	gormAccountRepo "github.com/EdlanioJ/kbu-store/infra/db/repository/gorm"
	gormStoreRepo "github.com/EdlanioJ/kbu-store/store/repository/gorm"
	"gorm.io/gorm"
)

func StoreUsecase(db *gorm.DB, contextTimeout time.Duration) *usecase.StoreUsecase {
	storeRepo := gormStoreRepo.NewStoreRepository(db)
	accountRepo := gormAccountRepo.NewAccountRepository(db)
	categoryRepo := gormCategoryRepo.NewCategoryRepository(db)

	return usecase.NewStoreUsecase(storeRepo, accountRepo, categoryRepo, contextTimeout)
}

func TagUsecase(db *gorm.DB, contextTimeout time.Duration) *usecase.TagUsecase {
	tagRepo := gormStoreRepo.NewTagsRepository(db)

	return usecase.NewtagUsecase(tagRepo, contextTimeout)
}
