package domain

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

const (
	// pending store status value
	StoreStatusPending string = "pending"
	// active store status value
	StoreStatusActive string = "active"
	// inactive store status value
	StoreStatusDisable string = "disable"
	// block store status value
	StoreStatusBlock string = "block"
)

// Stores belong to the domain layer.
type Stores []*Store

// A Store belong to the domain layer.
type Store struct {
	Base        `valid:"required"`
	Name        string                `json:"name" valid:"notnull"`
	Description string                `json:"description" valid:"-"`
	Status      string                `json:"status" valid:"notnull,status"`
	UserID      string                `json:"user_id" valid:"notnull,uuidv4"`
	AccountID   string                `json:"account_id" valid:"notnull,uuidv4"`
	Category    *Category             `valid:"-"`
	Image       *multipart.FileHeader `json:"-" valid:"-"`
	Image_Url   string                `json:"image_url" valid:"-"`
	Tags        []string              `json:"tags" valid:"-"`
	Position    Position              `json:"location" valid:"-"`
}

type (
	// StoreRepository represent the store's repository contract
	StoreRepository interface {
		Create(ctx context.Context, store *Store) error
		FindByID(ctx context.Context, id string) (*Store, error)
		FindByName(ctx context.Context, name string) (*Store, error)
		FindAll(ctx context.Context, sort string, limit, page int) (Stores, int64, error)
		Update(ctx context.Context, store *Store) error
		Delete(ctx context.Context, id string) error
	}

	// StoreUsecase represent the store's usecase contract
	StoreUsecase interface {
		Store(ctx context.Context, name, description, CategoryID, externalID string, tags []string, lat, lng float64) error
		Index(ctx context.Context, sort string, limit, page int) (Stores, int64, error)
		Get(ctx context.Context, id string) (*Store, error)
		Update(ctx context.Context, store *Store) error
		Delete(ctx context.Context, id string) error
		Block(ctx context.Context, id string) error
		Active(ctx context.Context, id string) error
		Disable(ctx context.Context, id string) error
	}
)

// Store entity validator
func (s *Store) isValid() (err error) {
	govalidator.TagMap["status"] = govalidator.Validator(func(str string) bool {
		return govalidator.IsIn(str, StoreStatusActive, StoreStatusPending, StoreStatusDisable, StoreStatusBlock)
	})

	_, err = govalidator.ValidateStruct(s)

	return
}

// Block set store entity status to block
func (s *Store) Block() (err error) {
	if s.Status == StoreStatusBlock {
		return ErrBlocked
	}

	if s.Status == StoreStatusPending {
		return ErrIsPending
	}

	s.Status = StoreStatusBlock

	err = s.isValid()
	return
}

// Activate set store entity status to active
func (s *Store) Activate() (err error) {
	if s.Status == StoreStatusActive {
		return ErrActived
	}

	s.Status = StoreStatusActive

	err = s.isValid()
	return
}

// Disable set store entity status to disable
func (s *Store) Disable() (err error) {
	if s.Status == StoreStatusDisable {
		return ErrInactived
	}

	if s.Status == StoreStatusBlock {
		return ErrBlocked
	}

	s.Status = StoreStatusDisable

	err = s.isValid()
	return
}

// NewStore creates a store entity
func NewStore(name, description, userID string, category *Category, accountID string, tags []string, lat, lng float64) (store *Store, err error) {
	store = &Store{
		Name:        name,
		Description: description,
		UserID:      userID,
		AccountID:   accountID,
		Category:    category,
		Tags:        tags,
	}

	store.ID = uuid.NewV4().String()
	store.Position = Position{
		Lat: lat,
		Lng: lng,
	}
	store.Status = StoreStatusPending
	store.CreatedAt = time.Now()

	err = store.isValid()

	if err != nil {
		return nil, err
	}
	return
}
