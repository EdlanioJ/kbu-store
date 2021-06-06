package factory

import (
	"time"

	gormAccountRepo "github.com/EdlanioJ/kbu-store/account/repository/gorm"
	gormCategoryRepo "github.com/EdlanioJ/kbu-store/category/repository/gorm"
	"github.com/EdlanioJ/kbu-store/store/handler"
	gormStoreRepo "github.com/EdlanioJ/kbu-store/store/repository/gorm"
	"github.com/EdlanioJ/kbu-store/store/usecase"
	"gorm.io/gorm"
)

func ActiveStoreHandler(db *gorm.DB, contextTimeout time.Duration) *handler.ActiveHandler {
	getStoreRepo := gormStoreRepo.NewGormGetStoreByID(db)
	updateStore := gormStoreRepo.NewGormUpdateStore(db)

	activeUsecase := usecase.NewActivateStore(getStoreRepo, updateStore, contextTimeout)
	return handler.NewActiveHandler(activeUsecase)
}

func BlockStoreHandler(db *gorm.DB, contextTimeout time.Duration) *handler.BlockHandler {
	getStoreRepo := gormStoreRepo.NewGormGetStoreByID(db)
	updateStore := gormStoreRepo.NewGormUpdateStore(db)

	blockUsecase := usecase.NewBlockStore(getStoreRepo, updateStore, contextTimeout)
	return handler.NewBlockHandler(blockUsecase)
}

func DisableStoreHandler(db *gorm.DB, contextTimeout time.Duration) *handler.DisableHandler {
	getStoreRepo := gormStoreRepo.NewGormGetStoreByID(db)
	updateStore := gormStoreRepo.NewGormUpdateStore(db)

	disableUsecase := usecase.NewDisableStore(getStoreRepo, updateStore, contextTimeout)
	return handler.NewDisableHandler(disableUsecase)
}

func CreateStoreHandler(db *gorm.DB, contextTimeout time.Duration) *handler.CreateHandler {
	createStoreRepo := gormStoreRepo.NewGormCreateStore(db)
	createAccountRepo := gormAccountRepo.NewGormCreateAccount(db)
	getCategoryRepo := gormCategoryRepo.NewGormGetCategoryByID(db)

	createUsecase := usecase.NewCreateStore(createStoreRepo, createAccountRepo, getCategoryRepo, contextTimeout)
	return handler.NewCreateHandler(createUsecase)
}

func DeleteStoreHandler(db *gorm.DB, contextTimeout time.Duration) *handler.DeleteHandler {
	deleteStoreRepo := gormStoreRepo.NewGormDeleteStore(db)
	deleteAccountRepo := gormAccountRepo.NewGormDeleteAccount(db)
	getStoreRepo := gormStoreRepo.NewGormGetStoreByID(db)

	deleteUsecase := usecase.NewDeleteStore(deleteStoreRepo, deleteAccountRepo, getStoreRepo, contextTimeout)
	return handler.NewDeleteHandler(deleteUsecase)
}

func FetchByLocationAndStatusHandler(db *gorm.DB, contextTimeout time.Duration) *handler.FetchByLocationAndStatusHandler {
	fetchByLocationRepo := gormStoreRepo.NewGormFetchStoreByCloseLocation(db)
	getCategoryRepo := gormCategoryRepo.NewGormGetCategoryByID(db)

	fetchByLocation := usecase.NewFetchStoreByCloseLocation(fetchByLocationRepo, getCategoryRepo, contextTimeout)
	return handler.NewFetchByLocationAndStatusHandler(fetchByLocation)
}
func FetchStoreHandler(db *gorm.DB, contextTimeout time.Duration) *handler.FetchHandler {
	fetchStoreRepo := gormStoreRepo.NewGormFetchStore(db)
	getCategoryRepo := gormCategoryRepo.NewGormGetCategoryByID(db)

	fetchUsecase := usecase.NewFetchStore(fetchStoreRepo, getCategoryRepo, contextTimeout)
	return handler.NewFetchHandler(fetchUsecase)
}

func FetchStoreByOwnerHandler(db *gorm.DB, contextTimeout time.Duration) *handler.FetchByOwnerHandler {
	fetchByOwnerRepo := gormStoreRepo.NewGormFetchStoreByOwner(db)
	getCategoryRepo := gormCategoryRepo.NewGormGetCategoryByID(db)

	fetchByOwnerUsecase := usecase.NewFetchStoreByOwner(fetchByOwnerRepo, getCategoryRepo, contextTimeout)
	return handler.NewFetchByOwnerHandler(fetchByOwnerUsecase)
}

func FetchStoreByStatusHandler(db *gorm.DB, contextTimeout time.Duration) *handler.FetchByStatusHandler {
	fetchByStatusRepo := gormStoreRepo.NewGormFetchStoreByStauts(db)
	getCategoryRepo := gormCategoryRepo.NewGormGetCategoryByID(db)

	fetchByStatusUsecase := usecase.NewFetchStoreByStatus(fetchByStatusRepo, getCategoryRepo, contextTimeout)
	return handler.NewFetchByStatusHandler(fetchByStatusUsecase)
}

func FetchStoreByTagsHandler(db *gorm.DB, contextTimeout time.Duration) *handler.FetchByTagsHandler {
	fetchByTagsRepo := gormStoreRepo.NewGormFetchStoreByTags(db)
	getCategoryRepo := gormCategoryRepo.NewGormGetCategoryByID(db)

	fetchByTagsUsecase := usecase.NewFetchStoreByTabs(fetchByTagsRepo, getCategoryRepo, contextTimeout)
	return handler.NewFetchByTagsHandler(fetchByTagsUsecase)
}

func FetchStoreByTypeHandler(db *gorm.DB, contextTimeout time.Duration) *handler.FetchByTypeHandler {
	fetchByTypeRepo := gormStoreRepo.NewGormFetchStoreByCategory(db)
	getCategoryRepo := gormCategoryRepo.NewGormGetCategoryByID(db)

	fetchByTypeUsecase := usecase.NewFetchStoreByCategory(fetchByTypeRepo, getCategoryRepo, contextTimeout)
	return handler.NewFetchByTypeHandler(fetchByTypeUsecase)
}

func GetStoreByIDHandler(db *gorm.DB, contextTimeout time.Duration) *handler.GetByIDHandler {
	getStoreRepo := gormStoreRepo.NewGormGetStoreByID(db)
	getStorTypeRepo := gormCategoryRepo.NewGormGetCategoryByID(db)

	getbyIDUsecase := usecase.NewGetStoreByID(getStoreRepo, getStorTypeRepo, contextTimeout)
	return handler.NewGetByIDHandler(getbyIDUsecase)
}

func GetStoreByOwnerHandler(db *gorm.DB, contextTimeout time.Duration) *handler.GetByOwnerHandler {
	getStoreRepo := gormStoreRepo.NewGormGetStoreByOwner(db)
	getStorTypeRepo := gormCategoryRepo.NewGormGetCategoryByID(db)

	getbyIDUsecase := usecase.NewGetStoreByOwner(getStoreRepo, getStorTypeRepo, contextTimeout)
	return handler.NewGetByOwnerHandler(getbyIDUsecase)
}

func UpdateStoreHandler(db *gorm.DB, contextTimeout time.Duration) *handler.UpdateHandler {
	updateStoreRepo := gormStoreRepo.NewGormUpdateStore(db)
	getStoreRepo := gormStoreRepo.NewGormGetStoreByID(db)

	updateUsecase := usecase.NewUpdateStore(getStoreRepo, updateStoreRepo, contextTimeout)
	return handler.NewUpdateHandler(updateUsecase)
}
