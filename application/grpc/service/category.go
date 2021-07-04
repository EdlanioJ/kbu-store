package service

import (
	"context"

	"github.com/EdlanioJ/kbu-store/application/grpc/pb"
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/golang/protobuf/ptypes/empty"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (s *categotyService) GetById(ctx context.Context, in *pb.Request) (*pb.Category, error) {
	_, err := uuid.FromString(in.Id)
	if err != nil {
		return nil, err
	}

	res, err := s.categoryUsecase.GetById(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Category{
		ID:        res.ID,
		Name:      res.Name,
		Status:    res.Status,
		CreatedAt: timestamppb.New(res.CreatedAt),
	}, nil
}
