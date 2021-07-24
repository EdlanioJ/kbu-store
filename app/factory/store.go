package factory

import (
	"time"

	"github.com/EdlanioJ/kbu-store/app/infrastructure/kafka"
	gormRepo "github.com/EdlanioJ/kbu-store/app/infrastructure/repository/gorm"
	"github.com/EdlanioJ/kbu-store/app/usecases"
	"gorm.io/gorm"
)

func StoreUsecase(db *gorm.DB, timeout time.Duration, brokers []string) *usecases.StoreUsecase {
	storeRepo := gormRepo.NewStoreRepository(db)
	accountRepo := gormRepo.NewAccountRepository(db)
	categoryRepo := gormRepo.NewCategoryRepository(db)
	kafkaProducer := kafka.NewKafkaProducer(brokers)

	return usecases.NewStoreUsecase(storeRepo, accountRepo, categoryRepo, kafkaProducer, timeout)
}
