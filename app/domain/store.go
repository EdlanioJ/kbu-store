package domain

import (
	"context"
	"encoding/json"
	"time"

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
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
	UserID      string   `json:"user_id"`
	AccountID   string   `json:"account_id"`
	CategoryID  string   `json:"category_id"`
	Image       string   `json:"image"`
	Tags        []string `json:"tags"`
	Position    Position `json:"location"`
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
func NewStore() (store *Store) {
	store = new(Store)

	store.ID = uuid.NewV4().String()
	store.Status = StoreStatusPending
	store.CreatedAt = time.Now()
	return
}
