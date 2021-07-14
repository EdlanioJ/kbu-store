package service

import (
	"context"

	"github.com/EdlanioJ/kbu-store/application/grpc/pb"
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/validators"
)

type tagService struct {
	tagUsecase domain.TagUsecase
	pb.UnimplementedTagServiceServer
}

func NewTagServer(usecase domain.TagUsecase) pb.TagServiceServer {
	return &tagService{
		tagUsecase: usecase,
	}
}

func (s *tagService) GetAll(ctx context.Context, in *pb.TagListRequest) (*pb.TagListResponse, error) {
	var tags []*pb.Tag
	res, total, err := s.tagUsecase.GetAll(ctx, in.Sort, int(in.Page), int(in.Limit))
	if err != nil {
		return nil, err
	}

	for _, item := range res {
		tags = append(tags, &pb.Tag{
			Tag:   item.Name,
			Count: int32(item.Count),
		})
	}

	return &pb.TagListResponse{
		Tags:  tags,
		Total: total,
	}, nil
}

func (s *tagService) GetAllByCategory(ctx context.Context, in *pb.TagListByCategoryRequest) (*pb.TagListResponse, error) {
	var tags []*pb.Tag
	err := validators.ValidateUUIDV4("category", in.CategoryId)
	if err != nil {
		return nil, err
	}

	res, total, err := s.tagUsecase.GetAllByCategory(ctx, in.CategoryId, in.Sort, int(in.Page), int(in.Limit))
	if err != nil {
		return nil, err
	}

	for _, item := range res {
		tags = append(tags, &pb.Tag{
			Tag:   item.Name,
			Count: int32(item.Count),
		})
	}

	return &pb.TagListResponse{
		Tags:  tags,
		Total: total,
	}, nil
}
