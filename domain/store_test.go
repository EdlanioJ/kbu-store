package domain_test

import (
	"testing"

	"github.com/EdlanioJ/kbu-store/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewStore(t *testing.T) {
	t.Run("should fail if validation fails", func(t *testing.T) {
		category, _ := domain.NewCategory("Type 001")
		account, _ := domain.NewAccount(200)
		tags := []string{"tag 001", "tag 002"}
		store, err := domain.NewStore("", "", "", category, account.ID, tags, 0, 0)

		assert.Nil(t, store)
		assert.Error(t, err)
	})

	t.Run("should succeed", func(t *testing.T) {
		name := "store 001"
		description := "store description 001"
		externalID := uuid.NewV4().String()
		category, _ := domain.NewCategory("Type 001")

		account, _ := domain.NewAccount(200)
		tags := []string{"tag 001", "tag 002"}
		store, err := domain.NewStore(name, description, externalID, category, account.ID, tags, 0, 0)

		assert.NotNil(t, store)
		assert.Nil(t, err)
	})
}

func TestBock(t *testing.T) {
	t.Run("should fail if status already bocked", func(t *testing.T) {
		name := "store 001"
		description := "store description 001"
		externalID := uuid.NewV4().String()
		category, _ := domain.NewCategory("Type 001")

		account, _ := domain.NewAccount(200)
		tags := []string{"tag 001", "tag 002"}
		store, err := domain.NewStore(name, description, externalID, category, account.ID, tags, 0, 0)
		assert.NoError(t, err)
		store.Status = domain.StoreStatusBlock
		store.ID = "001"
		err = store.Block()

		assert.Error(t, err)
	})

	t.Run("should fail if status is still pending", func(t *testing.T) {
		name := "store 001"
		description := "store description 001"
		externalID := uuid.NewV4().String()
		category, _ := domain.NewCategory("Type 001")

		account, _ := domain.NewAccount(200)
		tags := []string{"tag 001", "tag 002"}
		store, err := domain.NewStore(name, description, externalID, category, account.ID, tags, 0, 0)
		assert.NoError(t, err)
		store.ID = "001"
		err = store.Block()

		assert.Error(t, err)
	})

	t.Run("should fail if validation fails", func(t *testing.T) {
		name := "store 001"
		description := "store description 001"
		externalID := uuid.NewV4().String()
		category, _ := domain.NewCategory("Type 001")

		account, _ := domain.NewAccount(200)
		tags := []string{"tag 001", "tag 002"}
		store, err := domain.NewStore(name, description, externalID, category, account.ID, tags, 0, 0)
		assert.NoError(t, err)
		store.ID = "001"
		store.Status = domain.StoreStatusActive
		err = store.Block()

		assert.Error(t, err)
	})

	t.Run("should succeed", func(t *testing.T) {
		name := "store 001"
		description := "store description 001"
		externalID := uuid.NewV4().String()
		category, _ := domain.NewCategory("Type 001")
		account, _ := domain.NewAccount(200)
		tags := []string{"tag 001", "tag 002"}
		store, _ := domain.NewStore(name, description, externalID, category, account.ID, tags, 0, 0)
		store.Status = domain.StoreStatusActive
		err := store.Block()

		assert.Nil(t, err)
	})
}

func TestActivate(t *testing.T) {
	t.Run("should fail if status is already active", func(t *testing.T) {
		name := "store 001"
		description := "store description 001"
		externalID := uuid.NewV4().String()
		category, _ := domain.NewCategory("Type 001")

		account, _ := domain.NewAccount(200)
		tags := []string{"tag 001", "tag 002"}
		store, err := domain.NewStore(name, description, externalID, category, account.ID, tags, 0, 0)
		assert.NoError(t, err)
		store.Status = domain.StoreStatusActive
		err = store.Activate()

		assert.Error(t, err)
	})

	t.Run("should fail if validation fails", func(t *testing.T) {
		name := "store 001"
		description := "store description 001"
		externalID := uuid.NewV4().String()
		category, _ := domain.NewCategory("Type 001")

		account, _ := domain.NewAccount(200)
		tags := []string{"tag 001", "tag 002"}
		store, err := domain.NewStore(name, description, externalID, category, account.ID, tags, 0, 0)
		assert.NoError(t, err)

		store.ID = "001"
		err = store.Activate()

		assert.Error(t, err)
	})

	t.Run("should succeed", func(t *testing.T) {
		name := "store 001"
		description := "store description 001"
		externalID := uuid.NewV4().String()
		category, _ := domain.NewCategory("Type 001")
		account, _ := domain.NewAccount(200)
		tags := []string{"tag 001", "tag 002"}
		store, err := domain.NewStore(name, description, externalID, category, account.ID, tags, 0, 0)
		assert.NoError(t, err)
		err = store.Activate()

		assert.Nil(t, err)
	})
}

func TestDisable(t *testing.T) {
	t.Run("should fail if store status already disable", func(t *testing.T) {
		name := "store 001"
		description := "store description 001"
		externalID := uuid.NewV4().String()
		category, _ := domain.NewCategory("Type 001")

		account, _ := domain.NewAccount(200)
		tags := []string{"tag 001", "tag 002"}
		store, _ := domain.NewStore(name, description, externalID, category, account.ID, tags, 0, 0)
		store.Status = domain.StoreStatusDisable
		err := store.Disable()

		assert.Error(t, err)
	})

	t.Run("should fail if store status is blocked", func(t *testing.T) {
		name := "store 001"
		description := "store description 001"
		externalID := uuid.NewV4().String()
		category, _ := domain.NewCategory("Type 001")

		account, _ := domain.NewAccount(200)
		tags := []string{"tag 001", "tag 002"}
		store, _ := domain.NewStore(name, description, externalID, category, account.ID, tags, 0, 0)
		store.Status = domain.StoreStatusBlock
		err := store.Disable()

		assert.Error(t, err)
	})

	t.Run("should fail if validation fails", func(t *testing.T) {
		name := "store 001"
		description := "store description 001"
		externalID := uuid.NewV4().String()
		category, _ := domain.NewCategory("Type 001")

		account, _ := domain.NewAccount(200)
		tags := []string{"tag 001", "tag 002"}
		store, _ := domain.NewStore(name, description, externalID, category, account.ID, tags, 0, 0)

		store.ID = "001"
		err := store.Disable()

		assert.Error(t, err)
	})

	t.Run("should succeed", func(t *testing.T) {
		name := "store 001"
		description := "store description 001"
		externalID := uuid.NewV4().String()
		category, _ := domain.NewCategory("Type 001")
		account, _ := domain.NewAccount(200)
		tags := []string{"tag 001", "tag 002"}
		store, _ := domain.NewStore(name, description, externalID, category, account.ID, tags, 0, 0)

		err := store.Disable()

		assert.Nil(t, err)
	})
}
