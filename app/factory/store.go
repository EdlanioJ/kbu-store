package factory

import (
	"time"

	"github.com/EdlanioJ/kbu-store/app/config"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/kafka"
	gormRepo "github.com/EdlanioJ/kbu-store/app/infrastructure/repository/gorm"
	"github.com/EdlanioJ/kbu-store/app/usecases"
	"gorm.io/gorm"
)

func StoreUsecase(db *gorm.DB, timeout time.Duration, cfg *config.Config) *usecases.StoreUsecase {
	storeRepo := gormRepo.NewStoreRepository(db)
	accountRepo := gormRepo.NewAccountRepository(db)
	categoryRepo := gormRepo.NewCategoryRepository(db)
	kafkaProducer := kafka.NewKafkaProducer(cfg)

	return usecases.NewStoreUsecase(storeRepo, accountRepo, categoryRepo, kafkaProducer, timeout)
}
