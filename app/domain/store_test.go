package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/EdlanioJ/kbu-store/app/domain"
	"github.com/EdlanioJ/kbu-store/app/utils/sample"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Store(t *testing.T) {
	t.Parallel()

	cr := sample.NewCreateStoreRequest()
	ur := sample.NewUpdateStoreRequest()

	t.Run("new_store", func(t *testing.T) {
		store := domain.NewStore(cr)

		assert.NotNil(t, store)
		assert.NotEmpty(t, store.ID)
		assert.Equal(t, store.Status, domain.StoreStatusPending)
		assert.Equal(t, cr.CategoryID, store.CategoryID)
		assert.Equal(t, cr.Description, store.Description)
		assert.Equal(t, cr.Name, store.Name)
		assert.Equal(t, cr.Lat, store.Position.Lat)
		assert.Equal(t, cr.Lng, store.Position.Lng)
		assert.Equal(t, cr.UserID, store.UserID)
		assert.Equal(t, pq.StringArray(cr.Tags), store.Tags)
	})

	t.Run("from_update_request", func(t *testing.T) {
		store := new(domain.Store)
		store.FromUpdateRequest(ur)

		assert.NotNil(t, store)
		assert.Equal(t, store.ID, ur.ID)
		assert.Equal(t, store.CategoryID, ur.CategoryID)
		assert.Equal(t, store.Description, ur.Description)
		assert.Equal(t, store.Name, ur.Name)
		assert.Equal(t, store.Lat, ur.Lat)
		assert.Equal(t, store.Lng, ur.Lng)
		assert.Equal(t, store.Image, ur.Image)
		assert.Equal(t, store.Tags, pq.StringArray(ur.Tags))
	})

	t.Run("block", func(t *testing.T) {
		t.Parallel()
		t.Run("failure_already_blocked", func(t *testing.T) {
			store := domain.NewStore(cr)
			store.Status = domain.StoreStatusBlock
			err := store.Block()

			assert.Error(t, err)
		})

		t.Run("failure_still_pending", func(t *testing.T) {
			store := domain.NewStore(cr)
			err := store.Block()

			assert.Error(t, err)
		})

		t.Run("success", func(t *testing.T) {
			store := domain.NewStore(cr)
			store.Status = domain.StoreStatusActive
			err := store.Block()

			assert.Nil(t, err)
		})
	})

	t.Run("activate", func(t *testing.T) {
		t.Parallel()
		t.Run("failure_already_actived", func(t *testing.T) {
			store := domain.NewStore(cr)
			store.Status = domain.StoreStatusActive
			err := store.Activate()

			assert.Error(t, err)
		})

		t.Run("success", func(t *testing.T) {
			store := domain.NewStore(cr)
			err := store.Activate()

			assert.Nil(t, err)
		})
	})

	t.Run("disable", func(t *testing.T) {
		t.Parallel()
		t.Run("failure_already_disabled", func(t *testing.T) {
			store := domain.NewStore(cr)
			store.Status = domain.StoreStatusDisable
			err := store.Disable()

			assert.Error(t, err)
		})
		t.Run("failure_is_blocked", func(t *testing.T) {
			store := domain.NewStore(cr)
			store.Status = domain.StoreStatusBlock
			err := store.Disable()

			assert.Error(t, err)
		})
		t.Run("success", func(t *testing.T) {
			store := domain.NewStore(cr)
			err := store.Disable()

			assert.Nil(t, err)
		})
	})
	t.Run("to_json", func(t *testing.T) {
		store := domain.NewStore(cr)
		store.Name = "store 001"
		store.Description = "store description 001"
		store.AccountID = uuid.NewV4().String()
		store.CategoryID = uuid.NewV4().String()
		store.UserID = uuid.NewV4().String()
		store.Tags = []string{"tag 001", "tag 002"}

		data := store.ToJson()
		assert.NotNil(t, data)
		newStore := new(domain.Store)
		err := json.Unmarshal(data, newStore)
		assert.NoError(t, err)
		assert.Equal(t, store.ID, newStore.ID)
	})
}
