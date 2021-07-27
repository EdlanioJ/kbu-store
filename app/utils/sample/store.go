package sample

import (
	"time"

	"github.com/EdlanioJ/kbu-store/app/domain"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/grpc/pb"
	uuid "github.com/satori/go.uuid"
)

func NewStore() *domain.Store {
	store := &domain.Store{
		Base: domain.Base{
			ID:        uuid.NewV4().String(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        "Store 001",
		Description: "store description 001",
		Status:      domain.StoreStatusPending,
		UserID:      uuid.NewV4().String(),
		AccountID:   uuid.NewV4().String(),
		CategoryID:  uuid.NewV4().String(),
		Position: domain.Position{
			Lat: -8.8368200,
			Lng: 13.2343200,
		},
	}

	return store
}

func NewPBCreateStoreRequest() *pb.CreateStoreRequest {
	return &pb.CreateStoreRequest{
		Name:        "store 001",
		Description: "store description 001",
		CategoryID:  uuid.NewV4().String(),
		ExternalID:  uuid.NewV4().String(),
		Tags:        []string{"tag001", "tag002"},
		Latitude:    -8.8368200,
		Longitude:   13.2343200,
	}
}

func NewPBStoreRequest() *pb.StoreRequest {
	return &pb.StoreRequest{
		Id: uuid.NewV4().String(),
	}
}
func NewPBListStoreRequest() *pb.ListStoreRequest {
	return &pb.ListStoreRequest{
		Page:  1,
		Limit: 10,
		Sort:  "created_at",
	}
}

func NewPBUpdateStoreRequest() *pb.UpdateStoreRequest {
	return &pb.UpdateStoreRequest{
		ID:          uuid.NewV4().String(),
		Name:        "store 001",
		Description: "store description 001",
		CategoryID:  uuid.NewV4().String(),
		Tags:        []string{"tag001", "tag002"},
		Latitude:    -8.8368200,
		Longitude:   13.2343200,
	}
}
