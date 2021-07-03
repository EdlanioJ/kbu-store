package service

import (
	"context"

	"github.com/EdlanioJ/kbu-store/application/grpc/pb"
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/golang/protobuf/ptypes/empty"
)

type categotyService struct {
	categoryUsecase domain.CategoryUsecase
	pb.UnimplementedCategoryServiceServer
}

func NewCategotyServer(usecase domain.CategoryUsecase) pb.CategoryServiceServer {
	return &categotyService{
		categoryUsecase: usecase,
	}
}

func (s *categotyService) Create(ctx context.Context, in *pb.CreateRequest) (*empty.Empty, error) {
	err := s.categoryUsecase.Create(ctx, in.Name)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
