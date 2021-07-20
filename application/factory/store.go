package factory

import (
	"time"

	gormRepo "github.com/EdlanioJ/kbu-store/infra/db/repository/gorm"
	"github.com/EdlanioJ/kbu-store/infra/kafka"
	"github.com/EdlanioJ/kbu-store/usecase"
	"gorm.io/gorm"
)

func StoreUsecase(db *gorm.DB, timeout time.Duration, brokers []string) *usecase.StoreUsecase {
	storeRepo := gormRepo.NewStoreRepository(db)
	accountRepo := gormRepo.NewAccountRepository(db)
	categoryRepo := gormRepo.NewCategoryRepository(db)
	kafkaProducer := kafka.NewKafkaProducer(brokers)

	return usecase.NewStoreUsecase(storeRepo, accountRepo, categoryRepo, kafkaProducer, timeout)
}
