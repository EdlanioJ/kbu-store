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

func (s *categotyService) GetByIdAndStatus(ctx context.Context, in *pb.GetByIdAndStatusRequest) (*pb.Category, error) {
	_, err := uuid.FromString(in.Id)
	if err != nil {
		return nil, err
	}

	res, err := s.categoryUsecase.GetByIdAndStatus(ctx, in.Id, in.Status.String())
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

func (s *categotyService) GetAll(ctx context.Context, in *pb.GetAllRequest) (*pb.ListResponse, error) {
	var categories []*pb.Category
	res, total, err := s.categoryUsecase.GetAll(ctx, in.Sort, int(in.Page), int(in.Limit))
	if err != nil {
		return nil, err
	}

	for _, item := range res {
		categories = append(categories, &pb.Category{
			ID:        item.ID,
			Name:      item.Name,
			Status:    item.Status,
			CreatedAt: timestamppb.New(item.CreatedAt),
		})
	}
	return &pb.ListResponse{
		Categories: categories,
		Total:      total,
	}, nil
}

func (s *categotyService) GetAllByStatus(ctx context.Context, in *pb.GetAllByStatusRequest) (*pb.ListResponse, error) {
	var categories []*pb.Category
	res, total, err := s.categoryUsecase.GetAllByStatus(ctx, in.Status.String(), in.Sort, int(in.Page), int(in.Limit))
	if err != nil {
		return nil, err
	}

	for _, item := range res {
		categories = append(categories, &pb.Category{
			ID:        item.ID,
			Name:      item.Name,
			Status:    item.Status,
			CreatedAt: timestamppb.New(item.CreatedAt),
		})
	}
	return &pb.ListResponse{
		Categories: categories,
		Total:      total,
	}, nil
}
