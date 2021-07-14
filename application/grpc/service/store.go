package service

import (
	"context"

	"github.com/EdlanioJ/kbu-store/application/grpc/pb"
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/golang/protobuf/ptypes/empty"
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
