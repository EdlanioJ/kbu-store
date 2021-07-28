package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/EdlanioJ/kbu-store/app/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Store(t *testing.T) {
	t.Parallel()

	t.Run("new_store", func(t *testing.T) {
		store := domain.NewStore()

		assert.NotNil(t, store)
		assert.NotEmpty(t, store.ID)
		assert.Equal(t, store.Status, domain.StoreStatusPending)
	})

	t.Run("block", func(t *testing.T) {
		t.Parallel()
		t.Run("failure_already_blocked", func(t *testing.T) {
			store := domain.NewStore()
			store.Status = domain.StoreStatusBlock
			err := store.Block()

			assert.Error(t, err)
		})

		t.Run("failure_still_pending", func(t *testing.T) {
			store := domain.NewStore()
			err := store.Block()

			assert.Error(t, err)
		})

		t.Run("success", func(t *testing.T) {
			store := domain.NewStore()
			store.Status = domain.StoreStatusActive
			err := store.Block()

			assert.Nil(t, err)
		})
	})

	t.Run("activate", func(t *testing.T) {
		t.Parallel()
		t.Run("failure_already_actived", func(t *testing.T) {
			store := domain.NewStore()
			store.Status = domain.StoreStatusActive
			err := store.Activate()

			assert.Error(t, err)
		})

		t.Run("success", func(t *testing.T) {
			store := domain.NewStore()
			err := store.Activate()

			assert.Nil(t, err)
		})
	})

	t.Run("disable", func(t *testing.T) {
		t.Parallel()
		t.Run("failure_already_disabled", func(t *testing.T) {
			store := domain.NewStore()
			store.Status = domain.StoreStatusDisable
			err := store.Disable()

			assert.Error(t, err)
		})
		t.Run("failure_is_blocked", func(t *testing.T) {
			store := domain.NewStore()
			store.Status = domain.StoreStatusBlock
			err := store.Disable()

			assert.Error(t, err)
		})
		t.Run("success", func(t *testing.T) {
			store := domain.NewStore()
			err := store.Disable()

			assert.Nil(t, err)
		})
	})
	t.Run("to_json", func(t *testing.T) {
		store := domain.NewStore()
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
