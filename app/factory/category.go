package factory

import (
	"time"

	"github.com/EdlanioJ/kbu-store/category/handler"
	gormRepo "github.com/EdlanioJ/kbu-store/category/repository/gorm"
	"github.com/EdlanioJ/kbu-store/category/usecase"
	"gorm.io/gorm"
)

func CreateCategoryHandler(db *gorm.DB, contextTimeout time.Duration) *handler.CreateHandler {
	createRepo := gormRepo.NewGormCreateCategory(db)

	createUsecase := usecase.NewCreateCategory(createRepo, contextTimeout)
	return handler.NewCreateHandler(createUsecase)
}

func FetchCategoryHandler(db *gorm.DB, contextTimeout time.Duration) *handler.FetchHandler {
	fetchRepo := gormRepo.NewGormFetchCategory(db)

	fetchUsecase := usecase.NewFetchCategory(fetchRepo, contextTimeout)
	return handler.NewFetchHandler(fetchUsecase)
}

func FetchCategoryByStatusHandler(db *gorm.DB, contextTimeout time.Duration) *handler.FetchByStatusHandler {
	fetchRepo := gormRepo.NewGormFetchCategoryByStatus(db)

	fetchByStatusUsecase := usecase.NewFetchCategoryByStatus(fetchRepo, contextTimeout)
	return handler.NewFetchByStatusHandler(fetchByStatusUsecase)
}

func GetCategoryByIDHandler(db *gorm.DB, contextTimeout time.Duration) *handler.GetByIDHandler {
	getRepo := gormRepo.NewGormGetCategoryByID(db)

	getUsecase := usecase.NewGetCategoryByID(getRepo, contextTimeout)
	return handler.NewGetByIDHandler(getUsecase)
}

func GetCategoryByStatusHandler(db *gorm.DB, contextTimeout time.Duration) *handler.GetByStatusHandler {
	getRepo := gormRepo.NewGormGetCategoryByStatus(db)

	getUsecase := usecase.NewGetCategoryByStatus(getRepo, contextTimeout)
	return handler.NewGetByStatusHandler(getUsecase)
}
