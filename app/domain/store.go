package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/lib/pq"
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
	Base
	Name        string         `json:"name" gorm:"column:name;type:varchar;not null"`
	Description string         `json:"description" gorm:"type:varchar(255)"`
	Status      string         `json:"status" gorm:"type:varchar(20)"`
	UserID      string         `json:"user_id" gorm:"column:user_id;type:uuid"`
	AccountID   string         `json:"account_id" gorm:"column:account_id;type:uuid"`
	CategoryID  string         `json:"category_id" gorm:"column:category_id;type:uuid"`
	Image       string         `json:"image" gorm:"column:image;type:varchar(255)"`
	Tags        pq.StringArray `json:"tags" swaggertype:"array,string" gorm:"column:tags;type:text[]"`
	Position    `json:"location"`
}

type CreateStoreRequest struct {
	Name        string   `json:"name" validate:"required,alphanumunicode"`
	Description string   `json:"description" validate:"required,alphanumunicode"`
	CategoryID  string   `json:"category_id" validate:"required,uuid4"`
	UserID      string   `json:"user_id" validate:"required,uuid4"`
	Tags        []string `json:"tags"`
	Lat         float64  `json:"latitude" validate:"latitude"`
	Lng         float64  `json:"longitude" validate:"longitude"`
}

type UpdateStoreRequest struct {
	ID          string   `json:"-" validate:"required,uuid4"`
	Name        string   `json:"name" validate:"alphanumunicode"`
	Description string   `json:"description" validate:"alphanumunicode"`
	CategoryID  string   `json:"category_id" validate:"uuid4"`
	Image       string   `json:"image"`
	Tags        []string `json:"tags" validate:"uuid4"`
	Lat         float64  `json:"latitude" validate:"latitude"`
	Lng         float64  `json:"longitude" validate:"longitude"`
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
		Store(ctx context.Context, param *CreateStoreRequest) error
		Index(ctx context.Context, sort string, limit, page int) (Stores, int64, error)
		Get(ctx context.Context, id string) (*Store, error)
		Update(ctx context.Context, param *UpdateStoreRequest) error
		Delete(ctx context.Context, id string) error
		Block(ctx context.Context, id string) error
		Active(ctx context.Context, id string) error
		Disable(ctx context.Context, id string) error
	}
)

// Block set store entity status to block
func (s *Store) Block() (err error) {
	if s.Status == StoreStatusBlock {
		return ErrBlocked
	}

	if s.Status == StoreStatusPending {
		return ErrIsPending
	}

	s.Status = StoreStatusBlock
	return
}

// Activate set store entity status to active
func (s *Store) Activate() (err error) {
	if s.Status == StoreStatusActive {
		return ErrActived
	}

	s.Status = StoreStatusActive
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
	return
}

// ToJson returns the JSON encoding of Store
func (s *Store) ToJson() (res []byte) {
	res, _ = json.Marshal(s)

	return
}

// NewStore creates a store entity
func NewStore(param *CreateStoreRequest) (store *Store) {
	store = new(Store)

	store.ID = uuid.NewV4().String()
	store.Name = param.Name
	store.Description = param.Description
	store.CategoryID = param.CategoryID
	store.UserID = param.UserID
	store.Tags = param.Tags
	store.Position.Lat = param.Lat
	store.Position.Lng = param.Lng
	store.Status = StoreStatusPending
	store.CreatedAt = time.Now()
	return
}

func (s *Store) FromUpdateRequest(param *UpdateStoreRequest) {
	s.ID = param.ID
	s.Name = param.Name
	s.Description = param.Description
	s.CategoryID = param.CategoryID
	s.Image = param.Image
	s.Tags = param.Tags
	s.Position.Lat = param.Lat
	s.Position.Lng = param.Lng
	s.UpdatedAt = time.Now()
}
