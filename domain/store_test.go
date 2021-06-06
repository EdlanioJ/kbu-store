package domain_test

import (
	"testing"

	"github.com/EdlanioJ/kbu-store/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewStore(t *testing.T) {
	t.Run("should fail", func(t *testing.T) {
		is := assert.New(t)
		category, _ := domain.NewCategory("Type 001")
		account, _ := domain.NewAccount(200)
		tags := []string{"tag 001", "tag 002"}
		store, err := domain.NewStore("", "", "", category, account, tags, 0, 0)

		is.Nil(store)
		is.Error(err)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		name := "store 001"
		description := "store description 001"
		externalID := uuid.NewV4().String()
		category, _ := domain.NewCategory("Type 001")

		account, _ := domain.NewAccount(200)
		tags := []string{"tag 001", "tag 002"}
		store, err := domain.NewStore(name, description, externalID, category, account, tags, 0, 0)

		is.NotNil(store)
		is.Nil(err)
	})
}

func TestBock(t *testing.T) {
	t.Run("should fail", func(t *testing.T) {
		is := assert.New(t)

		name := "store 001"
		description := "store description 001"
		externalID := uuid.NewV4().String()
		category, _ := domain.NewCategory("Type 001")

		account, _ := domain.NewAccount(200)
		tags := []string{"tag 001", "tag 002"}
		store, _ := domain.NewStore(name, description, externalID, category, account, tags, 0, 0)

		store.ID = "001"
		err := store.Block()

		is.Error(err)
	})

	t.Run("should succeed", func(t *testing.T) {
		is := assert.New(t)

		name := "store 001"
		description := "store description 001"
		externalID := uuid.NewV4().String()
		category, _ := domain.NewCategory("Type 001")
		account, _ := domain.NewAccount(200)
		tags := []string{"tag 001", "tag 002"}
		store, _ := domain.NewStore(name, description, externalID, category, account, tags, 0, 0)

		err := store.Block()

		is.Nil(err)
	})
}

func TestActivate(t *testing.T) {
	t.Run("should fail", func(t *testing.T) {
		is := assert.New(t)

		name := "store 001"
		description := "store description 001"
		externalID := uuid.NewV4().String()
		category, _ := domain.NewCategory("Type 001")

		account, _ := domain.NewAccount(200)
		tags := []string{"tag 001", "tag 002"}
		store, _ := domain.NewStore(name, description, externalID, category, account, tags, 0, 0)

		store.ID = "001"
		err := store.Activate()

		is.Error(err)
	})

	t.Run("should succeed", func(t *testing.T) {
		is := assert.New(t)

		name := "store 001"
		description := "store description 001"
		externalID := uuid.NewV4().String()
		category, _ := domain.NewCategory("Type 001")
		account, _ := domain.NewAccount(200)
		tags := []string{"tag 001", "tag 002"}
		store, _ := domain.NewStore(name, description, externalID, category, account, tags, 0, 0)

		err := store.Activate()

		is.Nil(err)
	})
}

func TestInactivate(t *testing.T) {
	t.Run("should fail", func(t *testing.T) {
		is := assert.New(t)

		name := "store 001"
		description := "store description 001"
		externalID := uuid.NewV4().String()
		category, _ := domain.NewCategory("Type 001")

		account, _ := domain.NewAccount(200)
		tags := []string{"tag 001", "tag 002"}
		store, _ := domain.NewStore(name, description, externalID, category, account, tags, 0, 0)

		store.ID = "001"
		err := store.Inactivate()

		is.Error(err)
	})

	t.Run("should succeed", func(t *testing.T) {
		is := assert.New(t)

		name := "store 001"
		description := "store description 001"
		externalID := uuid.NewV4().String()
		category, _ := domain.NewCategory("Type 001")
		account, _ := domain.NewAccount(200)
		tags := []string{"tag 001", "tag 002"}
		store, _ := domain.NewStore(name, description, externalID, category, account, tags, 0, 0)

		err := store.Inactivate()

		is.Nil(err)
	})
}
