package service

import (
	"context"

	"github.com/EdlanioJ/kbu-store/application/grpc/pb"
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/validators"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type storeService struct {
	storeUsecase domain.StoreUsecase
	pb.UnimplementedStoreServiceServer
}

func NewStoreServer(u domain.StoreUsecase) pb.StoreServiceServer {
	return &storeService{
		storeUsecase: u,
	}
}

func (s *storeService) newPBStore(store *domain.Store) *pb.Store {
	t := &pb.Store{
		ID:          store.ID,
		Name:        store.Name,
		Description: store.Description,
		Status:      store.Status,
		ExternalID:  store.ExternalID,
		AccountID:   store.AccountID,
		Tags:        store.Tags,
		Location: &pb.Location{
			Latitude:  store.Position.Lat,
			Longitude: store.Position.Lng,
		},
		Category: &pb.Category{
			ID:        store.Category.ID,
			Name:      store.Category.Name,
			Status:    store.Category.Status,
			CreatedAt: timestamppb.New(store.Category.CreatedAt),
		},
		CreatedAt: timestamppb.New(store.CreatedAt),
	}
	return t
}

func (s *storeService) Create(ctx context.Context, in *pb.CreateStoreRequest) (*empty.Empty, error) {
	err := s.storeUsecase.Create(
		ctx,
		in.GetName(),
		in.GetDescription(),
		in.GetCategoryID(),
		in.GetExternalID(),
		in.GetTags(),
		in.GetLatitude(),
		in.GetLongitude(),
	)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *storeService) GetById(ctx context.Context, in *pb.StoreRequest) (*pb.Store, error) {
	err := validators.ValidateUUIDV4("id", in.GetId())
	if err != nil {
		return nil, err
	}
	res, err := s.storeUsecase.GetById(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	store := s.newPBStore(res)
	return store, nil
}

func (s *storeService) GetByIdAndOwner(ctx context.Context, in *pb.GetStoreByIdAndOwnerRequest) (*pb.Store, error) {
	err := validators.ValidateUUIDV4("id", in.GetID())
	if err != nil {
		return nil, err
	}

	err = validators.ValidateUUIDV4("id", in.GetOwner())
	if err != nil {
		return nil, err
	}

	res, err := s.storeUsecase.GetByIdAndOwner(ctx, in.GetID(), in.GetOwner())
	if err != nil {
		return nil, err
	}

	store := s.newPBStore(res)
	return store, nil
}

func (s *storeService) GetAll(ctx context.Context, in *pb.GetAllStoreRequest) (*pb.ListStoreResponse, error) {
	var stores []*pb.Store

	res, total, err := s.storeUsecase.GetAll(ctx, in.GetSort(), int(in.GetLimit()), int(in.GetPage()))
	if err != nil {
		return nil, err
	}
	for _, item := range res {
		stores = append(stores, s.newPBStore(item))
	}

	return &pb.ListStoreResponse{
		Stores: stores,
		Total:  total,
	}, nil
}

func (s *storeService) GetAllByCategory(ctx context.Context, in *pb.ListStoreRequest) (*pb.ListStoreResponse, error) {
	var stores []*pb.Store
	err := validators.ValidateUUIDV4("id", in.GetId())
	if err != nil {
		return nil, err
	}

	res, total, err := s.storeUsecase.GetAllByCategory(ctx, in.GetId(), in.GetSort(), int(in.GetLimit()), int(in.GetPage()))
	if err != nil {
		return nil, err
	}

	for _, item := range res {
		stores = append(stores, s.newPBStore(item))
	}

	return &pb.ListStoreResponse{
		Stores: stores,
		Total:  total,
	}, nil
}
