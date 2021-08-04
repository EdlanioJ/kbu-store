package service

import (
	"context"

	"github.com/EdlanioJ/kbu-store/app/domain"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/grpc/pb"
	"github.com/EdlanioJ/kbu-store/app/validators"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/opentracing/opentracing-go"
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
		ExternalID:  store.UserID,
		AccountID:   store.AccountID,
		Tags:        store.Tags,
		Location: &pb.Location{
			Latitude:  store.Position.Lat,
			Longitude: store.Position.Lng,
		},
		Category:  store.CategoryID,
		CreatedAt: timestamppb.New(store.CreatedAt),
	}
	return t
}

func (s *storeService) Create(ctx context.Context, in *pb.CreateStoreRequest) (*empty.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storeService.Create")
	defer span.Finish()

	err := s.storeUsecase.Store(
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

func (s *storeService) Get(ctx context.Context, in *pb.StoreRequest) (*pb.Store, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storeService.Get")
	defer span.Finish()

	err := validators.ValidateRequired("id", in.GetId())
	if err != nil {
		return nil, err
	}
	err = validators.ValidateUUIDV4("id", in.GetId())
	if err != nil {
		return nil, err
	}
	res, err := s.storeUsecase.Get(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	store := s.newPBStore(res)
	return store, nil
}

func (s *storeService) List(ctx context.Context, in *pb.ListStoreRequest) (*pb.ListStoreResponse, error) {
	var stores []*pb.Store

	span, ctx := opentracing.StartSpanFromContext(ctx, "storeService.List")
	defer span.Finish()

	res, total, err := s.storeUsecase.Index(ctx, in.GetSort(), int(in.GetLimit()), int(in.GetPage()))
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

func (s *storeService) Activate(ctx context.Context, in *pb.StoreRequest) (*empty.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storeService.Activate")
	defer span.Finish()

	err := validators.ValidateRequired("id", in.GetId())
	if err != nil {
		return nil, err
	}
	err = validators.ValidateUUIDV4("id", in.GetId())
	if err != nil {
		return nil, err
	}

	err = s.storeUsecase.Active(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *storeService) Block(ctx context.Context, in *pb.StoreRequest) (*empty.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storeService.Block")
	defer span.Finish()

	err := validators.ValidateRequired("id", in.GetId())
	if err != nil {
		return nil, err
	}
	err = validators.ValidateUUIDV4("id", in.GetId())
	if err != nil {
		return nil, err
	}

	err = s.storeUsecase.Block(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *storeService) Disable(ctx context.Context, in *pb.StoreRequest) (*empty.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storeService.Disable")
	defer span.Finish()

	err := validators.ValidateRequired("id", in.GetId())
	if err != nil {
		return nil, err
	}
	err = validators.ValidateUUIDV4("id", in.GetId())
	if err != nil {
		return nil, err
	}

	err = s.storeUsecase.Disable(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *storeService) Update(ctx context.Context, in *pb.UpdateStoreRequest) (*empty.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storeService.Update")
	defer span.Finish()

	err := validators.ValidateRequired("id", in.GetID())
	if err != nil {
		return nil, err
	}
	err = validators.ValidateUUIDV4("id", in.GetID())
	if err != nil {
		return nil, err
	}

	err = validators.ValidateUUIDV4("category", in.GetCategoryID())
	if err != nil {
		return nil, err
	}
	err = validators.ValidateLatitude(in.GetLatitude())
	if err != nil {
		return nil, err
	}

	err = validators.ValidateLongitude(in.GetLongitude())
	if err != nil {
		return nil, err
	}

	store := new(domain.Store)
	store.ID = in.GetID()
	store.Name = in.GetName()
	store.Description = in.GetDescription()
	store.Tags = in.GetTags()
	store.Position.Lat = in.GetLatitude()
	store.Position.Lng = in.GetLongitude()
	store.CategoryID = in.GetCategoryID()

	err = s.storeUsecase.Update(ctx, store)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *storeService) Delete(ctx context.Context, in *pb.StoreRequest) (*empty.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storeService.Delete")
	defer span.Finish()

	err := validators.ValidateRequired("id", in.GetId())
	if err != nil {
		return nil, err
	}
	err = validators.ValidateUUIDV4("id", in.GetId())
	if err != nil {
		return nil, err
	}

	err = s.storeUsecase.Delete(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
