package service

import (
	"context"

	"github.com/EdlanioJ/kbu-store/app/domain"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/grpc/pb"
	"github.com/go-playground/validator/v10"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type storeService struct {
	storeUsecase domain.StoreUsecase
	validate     *validator.Validate
	pb.UnimplementedStoreServiceServer
}

func NewStoreServer(usecase domain.StoreUsecase, validate *validator.Validate) pb.StoreServiceServer {
	return &storeService{
		storeUsecase: usecase,
		validate:     validate,
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

	cr := new(domain.CreateStoreRequest)
	cr.Name = in.GetName()
	cr.Description = in.GetDescription()
	cr.CategoryID = in.GetCategoryID()
	cr.UserID = in.GetExternalID()
	cr.Tags = in.GetTags()
	cr.Lat = in.GetLatitude()
	cr.Lng = in.GetLongitude()

	if err := s.validate.StructCtx(ctx, cr); err != nil {
		log.
			WithContext(ctx).
			Errorf("validate.StructCtx: %v", err)
		return nil, err
	}

	err := s.storeUsecase.Store(ctx, cr)
	if err != nil {
		log.
			WithContext(ctx).
			Errorf("storeUsecase.Store: %v", err)
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *storeService) Get(ctx context.Context, in *pb.StoreRequest) (*pb.Store, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storeService.Get")
	defer span.Finish()

	if err := s.validate.VarCtx(ctx, in.GetId(), "uuid4"); err != nil {
		log.
			WithContext(ctx).
			Errorf("validate.VarCtx: %v", err)
		return nil, err
	}

	res, err := s.storeUsecase.Get(ctx, in.GetId())
	if err != nil {
		log.
			WithContext(ctx).
			Errorf("storeUsecase.Get: %v", err)
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
		log.
			WithContext(ctx).
			Errorf("storeUsecase.Index: %v", err)
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

	if err := s.validate.VarCtx(ctx, in.GetId(), "uuid4"); err != nil {
		log.
			WithContext(ctx).
			Errorf("validate.VarCtx: %v", err)
		return nil, err
	}

	err := s.storeUsecase.Active(ctx, in.GetId())
	if err != nil {
		log.
			WithContext(ctx).
			Errorf("storeUsecase.Active: %v", err)
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *storeService) Block(ctx context.Context, in *pb.StoreRequest) (*empty.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storeService.Block")
	defer span.Finish()

	if err := s.validate.VarCtx(ctx, in.GetId(), "uuid4"); err != nil {
		log.
			WithContext(ctx).
			Errorf("validate.VarCtx: %v", err)
		return nil, err
	}

	err := s.storeUsecase.Block(ctx, in.GetId())
	if err != nil {
		log.
			WithContext(ctx).
			Errorf("storeUsecase.Block: %v", err)
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *storeService) Disable(ctx context.Context, in *pb.StoreRequest) (*empty.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storeService.Disable")
	defer span.Finish()

	if err := s.validate.VarCtx(ctx, in.GetId(), "uuid4"); err != nil {
		log.
			WithContext(ctx).
			Errorf("validate.VarCtx: %v", err)
		return nil, err
	}

	err := s.storeUsecase.Disable(ctx, in.GetId())
	if err != nil {
		log.
			WithContext(ctx).
			Errorf("storeUsecase.Disable: %v", err)
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *storeService) Update(ctx context.Context, in *pb.UpdateStoreRequest) (*empty.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storeService.Update")
	defer span.Finish()

	ur := new(domain.UpdateStoreRequest)

	ur.ID = in.GetID()
	ur.Name = in.GetName()
	ur.Description = in.GetDescription()
	ur.Tags = in.GetTags()
	ur.Lat = in.GetLatitude()
	ur.Image = in.GetImage()
	ur.Lng = in.GetLongitude()
	ur.CategoryID = in.GetCategoryID()

	if err := s.validate.StructCtx(ctx, ur); err != nil {
		log.
			WithContext(ctx).
			Errorf("validate.StructCtx: %v", err)
		return nil, err
	}
	err := s.storeUsecase.Update(ctx, ur)
	if err != nil {
		log.
			WithContext(ctx).
			Errorf("storeUsecase.Update: %v", err)
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *storeService) Delete(ctx context.Context, in *pb.StoreRequest) (*empty.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storeService.Delete")
	defer span.Finish()

	if err := s.validate.VarCtx(ctx, in.GetId(), "uuid4"); err != nil {
		log.
			WithContext(ctx).
			Errorf("validate.VarCtx: %v", err)
		return nil, err
	}

	err := s.storeUsecase.Delete(ctx, in.GetId())
	if err != nil {
		log.
			WithContext(ctx).
			Errorf("storeUsecase.Delete: %v", err)
		return nil, err
	}

	return &empty.Empty{}, nil
}
