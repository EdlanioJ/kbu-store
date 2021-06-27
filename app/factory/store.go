package factory

import (
	"time"

	gormAccountRepo "github.com/EdlanioJ/kbu-store/account/repository/gorm"
	gormCategoryRepo "github.com/EdlanioJ/kbu-store/category/repository/gorm"
	dataUsecase "github.com/EdlanioJ/kbu-store/data/usecase"
	gormStoreRepo "github.com/EdlanioJ/kbu-store/store/repository/gorm"
	"github.com/EdlanioJ/kbu-store/store/usecase"
	"gorm.io/gorm"
)

func StoreUsecase(db *gorm.DB, contextTimeout time.Duration) *dataUsecase.StoreUsecase {
	storeRepo := gormStoreRepo.NewStoreRepository(db)
	accountRepo := gormAccountRepo.NewAccountRepository(db)
	categoryRepo := gormCategoryRepo.NewCategoryRepository(db)

	return dataUsecase.NewStoreUsecase(storeRepo, accountRepo, categoryRepo, contextTimeout)
}

func TagUsecase(db *gorm.DB, contextTimeout time.Duration) *usecase.TagUsecase {
	tagRepo := gormStoreRepo.NewTagsRepository(db)

	return usecase.NewtagUsecase(tagRepo, contextTimeout)
}
