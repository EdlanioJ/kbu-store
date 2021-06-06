package domain

import (
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

const (
	// pendding store status value
	StoreStatusPending string = "pending"
	// active store status value
	StoreStatusActive string = "active"
	// inactive store status value
	StoreStatusInactive string = "disable"
	// block store status value
	StoreStatusBlock string = "block"
)

// Store struct
type Store struct {
	Base        `valid:"required"`
	Name        string    `json:"name" valid:"notnull"`
	Description string    `json:"description" valid:"-"`
	Status      string    `json:"status" valid:"notnull,status"`
	ExternalID  string    `json:"external_id" valid:"notnull,uuidv4"`
	Account     *Account  `valid:"-"`
	Category    *Category `valid:"-"`
	Tags        []string  `json:"tags" valid:"-"`
	Position    Position  `json:"location" valid:"-"`
}

// Repositories
type (
	FetchStoreRepository interface {
		Exec(ctx context.Context, sort string, limit, page int) ([]*Store, int64, error)
	}
	FetchStoreByStatusRepository interface {
		Exec(ctx context.Context, status, sort string, limit, page int) ([]*Store, int64, error)
	}
	FetchStoreByTypeRepository interface {
		Exec(ctx context.Context, typeID, sort string, limit, page int) ([]*Store, int64, error)
	}
	FetchStoreByOwnerRepository interface {
		Exec(ctx context.Context, externalID, sort string, limit, page int) ([]*Store, int64, error)
	}
	FetchStoreByTagsRepository interface {
		Exec(ctx context.Context, tags []string, sort string, limit, page int) ([]*Store, int64, error)
	}
	FetchStoreByCloseLocationRepository interface {
		Exec(ctx context.Context, lat, lng float64, distance, limit, page int, status, sort string) ([]*Store, int64, error)
	}
	GetStoreByIDRepository interface {
		Exec(ctx context.Context, id string) (*Store, error)
	}
	GetStoreByOwnerRepository interface {
		Exec(ctx context.Context, id string, externalID string) (*Store, error)
	}
	CreateStoreRepository interface {
		Add(ctx context.Context, store *Store) error
	}
	UpdateStoreRepository interface {
		Exec(ctx context.Context, store *Store) error
	}
	DeleteStoreRepository interface {
		Exec(ctx context.Context, id string) error
	}
)

// Usecases
type (
	// FetchStoreUsecase represents get all store's usecase
	FetchStoreUsecase interface {
		Exec(ctx context.Context, sort string, limit, page int) ([]*Store, int64, error)
	}
	// FetchStoreByStatusUsecase represents get all by status store's usecase
	FetchStoreByStatusUsecase interface {
		Exec(ctx context.Context, status, sort string, limit, page int) ([]*Store, int64, error)
	}
	// FetchStoreByTypeUsecase represents get all by type id store's usecase
	FetchStoreByTypeUsecase interface {
		Exec(ctx context.Context, typeID, sort string, limit, page int) ([]*Store, int64, error)
	}
	// FetchStoreByTagsUsecase represents get all by tags store's usecase
	FetchStoreByTagsUsecase interface {
		Exec(ctx context.Context, tags []string, sort string, limit, page int) ([]*Store, int64, error)
	}
	// FetchStoreByOwnerUsecase represents get all by owner store's usecase
	FetchStoreByOwnerUsecase interface {
		Exec(ctx context.Context, owner, sort string, limit, page int) ([]*Store, int64, error)
	}
	// GetStoreByIDUsecase represents get by id store's usecase
	GetStoreByIDUsecase interface {
		Exec(ctx context.Context, id string) (*Store, error)
	}
	// FetchStoreByCloseLocationUsecase represents get all by close location and status store's usecase
	FetchStoreByCloseLocationUsecase interface {
		Exec(ctx context.Context, lat, lng float64, distance int, status string, limit, page int, sort string) ([]*Store, int64, error)
	}
	// GetStoreByOwnerUsecase represents get one by owner store's usecase
	GetStoreByOwnerUsecase interface {
		Exec(ctx context.Context, id string, externalID string) (*Store, error)
	}
	// BlockStoreUsecase represents block store's usecase
	BlockStoreUsecase interface {
		Exec(ctx context.Context, id string) error
	}
	// ActivateStoreUsecase represents activate store's usecase
	ActivateStoreUsecase interface {
		Exec(ctx context.Context, id string) error
	}
	// DisableStoreUsecase represents inactivate store's usecase
	DisableStoreUsecase interface {
		Exec(ctx context.Context, id string) error
	}
	// CreateStoreUsecase represents create store's usecase
	CreateStoreUsecase interface {
		Add(ctx context.Context, name, description, CategoryID, externalID string, tags []string, lat, lng float64) error
	}
	// UpdateStoreUsecase represents update store's usecase
	UpdateStoreUsecase interface {
		Exec(ctx context.Context, store *Store) error
	}
	// DeleteStoreUsecase represents delete store's usecase
	DeleteStoreUsecase interface {
		Exec(ctx context.Context, id string) error
	}
)

// Store entity validator
func (s *Store) isValid() (err error) {
	govalidator.TagMap["status"] = govalidator.Validator(func(str string) bool {
		return govalidator.IsIn(str, StoreStatusActive, StoreStatusPending, StoreStatusInactive, StoreStatusBlock)
	})

	_, err = govalidator.ValidateStruct(s)

	return
}

func (s *Store) Block() (err error) {
	s.Status = StoreStatusBlock

	err = s.isValid()
	return
}

func (s *Store) Activate() (err error) {
	s.Status = StoreStatusActive

	err = s.isValid()
	return
}

func (s *Store) Inactivate() (err error) {
	s.Status = StoreStatusInactive

	err = s.isValid()
	return
}

// NewStore creates a *Store struct
func NewStore(name, description, externalID string, category *Category, account *Account, tags []string, lat, lng float64) (store *Store, err error) {
	store = &Store{
		Name:        name,
		Description: description,
		ExternalID:  externalID,
		Account:     account,
		Category:    category,
		Tags:        tags,
	}

	store.Position = Position{
		Lat: lat,
		Lng: lng,
	}
	store.Status = StoreStatusPending
	store.ID = uuid.NewV4().String()
	store.CreatedAt = time.Now()

	err = store.isValid()

	if err != nil {
		return nil, err
	}

	return
}
