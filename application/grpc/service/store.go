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

	return &pb.Store{
		ID:          res.ID,
		Name:        res.Name,
		Description: res.Description,
		Status:      res.Status,
		ExternalID:  res.ExternalID,
		AccountID:   res.AccountID,
		Tags:        res.Tags,
		Location: &pb.Location{
			Latitude:  res.Position.Lat,
			Longitude: res.Position.Lng,
		},
		Category: &pb.Category{
			ID:        res.Category.ID,
			Name:      res.Category.Name,
			Status:    res.Category.Status,
			CreatedAt: timestamppb.New(res.Category.CreatedAt),
		},
		CreatedAt: timestamppb.New(res.CreatedAt),
	}, nil
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

	return &pb.Store{
		ID:          res.ID,
		Name:        res.Name,
		Description: res.Description,
		Status:      res.Status,
		ExternalID:  res.ExternalID,
		AccountID:   res.AccountID,
		Tags:        res.Tags,
		Location: &pb.Location{
			Latitude:  res.Position.Lat,
			Longitude: res.Position.Lng,
		},
		Category: &pb.Category{
			ID:        res.Category.ID,
			Name:      res.Category.Name,
			Status:    res.Category.Status,
			CreatedAt: timestamppb.New(res.Category.CreatedAt),
		},
		CreatedAt: timestamppb.New(res.CreatedAt),
	}, nil
}

func (s *storeService) GetAll(ctx context.Context, in *pb.GetAllStoreRequest) (*pb.ListStoreResponse, error) {
	var stores []*pb.Store

	res, total, err := s.storeUsecase.GetAll(ctx, in.GetSort(), int(in.GetLimit()), int(in.GetPage()))
	if err != nil {
		return nil, err
	}
	for _, item := range res {
		stores = append(stores, &pb.Store{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Status:      item.Status,
			ExternalID:  item.ExternalID,
			AccountID:   item.AccountID,
			Tags:        item.Tags,
			Location: &pb.Location{
				Latitude:  item.Position.Lat,
				Longitude: item.Position.Lng,
			},
			Category: &pb.Category{
				ID:        item.Category.ID,
				Name:      item.Category.Name,
				Status:    item.Category.Status,
				CreatedAt: timestamppb.New(item.Category.CreatedAt),
			},
			CreatedAt: timestamppb.New(item.CreatedAt),
		})
	}

	return &pb.ListStoreResponse{
		Stores: stores,
		Total:  total,
	}, nil
}
