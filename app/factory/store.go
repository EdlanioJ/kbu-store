package factory

import (
	"time"

	gormAccountRepo "github.com/EdlanioJ/kbu-store/account/repository/gorm"
	gormCategoryRepo "github.com/EdlanioJ/kbu-store/category/repository/gorm"
	gormStoreRepo "github.com/EdlanioJ/kbu-store/store/repository/gorm"
	"github.com/EdlanioJ/kbu-store/store/usecase"
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
