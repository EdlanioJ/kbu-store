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
	AccountID   string    `json:"account_id" valid:"notnull,uuidv4"`
	Category    *Category `valid:"-"`
	Tags        []string  `json:"tags" valid:"-"`
	Position    Position  `json:"location" valid:"-"`
}

// Repositories
type (
	StoreRepository interface {
		Create(ctx context.Context, store *Store) error
		GetById(ctx context.Context, id string) (*Store, error)
		GetByIdAndOwner(ctx context.Context, id string, externalID string) (*Store, error)
		GetAll(ctx context.Context, sort string, limit, page int) ([]*Store, int64, error)
		GetAllByCategory(ctx context.Context, categoryID, sort string, limit, page int) ([]*Store, int64, error)
		GetAllByLocation(ctx context.Context, lat, lng float64, distance, limit, page int, status, sort string) ([]*Store, int64, error)
		GetAllByOwner(ctx context.Context, externalID, sort string, limit, page int) ([]*Store, int64, error)
		GetAllByTags(ctx context.Context, tags []string, sort string, limit, page int) ([]*Store, int64, error)
		GetAllByStatus(ctx context.Context, status, sort string, limit, page int) ([]*Store, int64, error)
		Update(ctx context.Context, store *Store) error
		Delete(ctx context.Context, id string) error
	}

	StoreUsecase interface {
		Create(ctx context.Context, name, description, CategoryID, externalID string, tags []string, lat, lng float64) error
		GetById(ctx context.Context, id string) (*Store, error)
		GetByIdAndOwner(ctx context.Context, id string, externalID string) (*Store, error)
		GetAll(ctx context.Context, sort string, limit, page int) ([]*Store, int64, error)
		GetAllByCategory(ctx context.Context, categoryID, sort string, limit, page int) ([]*Store, int64, error)
		GetAllByOwner(ctx context.Context, owner, sort string, limit, page int) ([]*Store, int64, error)
		GetAllByStatus(ctx context.Context, status, sort string, limit, page int) ([]*Store, int64, error)
		GetAllByTags(ctx context.Context, tags []string, sort string, limit, page int) ([]*Store, int64, error)
		GetAllByByCloseLocation(ctx context.Context, lat, lng float64, distance int, status string, limit, page int, sort string) ([]*Store, int64, error)
		Block(ctx context.Context, id string) error
		Active(ctx context.Context, id string) error
		Disable(ctx context.Context, id string) error
		Update(ctx context.Context, store *Store) error
		Delete(ctx context.Context, id string) error
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
func NewStore(name, description, externalID string, category *Category, accountID string, tags []string, lat, lng float64) (store *Store, err error) {
	store = &Store{
		Name:        name,
		Description: description,
		ExternalID:  externalID,
		AccountID:   accountID,
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
