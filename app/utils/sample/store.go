package sample

import (
	"time"

	"github.com/EdlanioJ/kbu-store/app/domain"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/grpc/pb"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/http/handler"
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
func NewHttpCreateStoreRequest() handler.CreateStoreRequest {
	return handler.CreateStoreRequest{
		Name:        "Store 001",
		Description: "Store description",
		CategoryID:  uuid.NewV4().String(),
		UserID:      uuid.NewV4().String(),
		Tags:        []string{"tag001", "tag002"},
		Lat:         -8.8867698,
		Lng:         13.4771186,
	}
}

func NewHttpUpdateStoreRequest() handler.UpdateStoreRequest {
	return handler.UpdateStoreRequest{
		Name:        "store 002",
		Description: "description 002",
		CategoryID:  uuid.NewV4().String(),
		Tags:        []string{"tag002", "tag003"},
		Lat:         -8.8867698,
		Lng:         13.4771186,
	}
}

type HttpListRequest struct {
	Sort  string
	Page  int
	Limit int
}

func NewHttpListReq() HttpListRequest {
	return HttpListRequest{
		Sort:  "created_at",
		Page:  1,
		Limit: 5,
	}
}

type StoreUsecaseCreateRequest struct {
	Name        string
	Description string
	CategoryID  string
	UserID      string
	Tags        []string
}

func NewStoreUsecaseRequest() StoreUsecaseCreateRequest {
	return StoreUsecaseCreateRequest{
		Name:        "Store 001",
		Description: "Store description 001",
		CategoryID:  uuid.NewV4().String(),
		UserID:      uuid.NewV4().String(),
		Tags:        []string{"tag002", "tag003"},
	}
}
