package domain_test

import (
	"fmt"
	"testing"

	"github.com/EdlanioJ/kbu-store/app/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewCategory(t *testing.T) {
	t.Run("should fail on validation", func(t *testing.T) {
		is := assert.New(t)

		category, err := domain.NewCategory("", "", "")

		is.Nil(category)
		is.Error(err)
	})

	t.Run("should succeed", func(t *testing.T) {
		is := assert.New(t)

		id := uuid.NewV4().String()
		name := "Type 001"
		category, err := domain.NewCategory(id, name, domain.CategoryStatusInactive)

		is.NotNil(category)
		is.Nil(err)
		is.Equal(category.Name, name)
	})
}

func Test_DomainCategory_Activate(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		id := uuid.NewV4().String()
		name := "Type 001"
		category, _ := domain.NewCategory(id, name, domain.CategoryStatusInactive)
		category.ID = ""
		err := category.Activate()
		assert.Error(t, err)
	})
	t.Run("success", func(t *testing.T) {
		id := uuid.NewV4().String()
		name := "Type 001"
		category, _ := domain.NewCategory(id, name, domain.CategoryStatusInactive)

		err := category.Activate()
		assert.NoError(t, err)
	})
}

func Test_DomainCategory_Disable(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		id := uuid.NewV4().String()
		name := "Type 001"
		category, _ := domain.NewCategory(id, name, domain.CategoryStatusInactive)
		category.ID = ""
		err := category.Disable()
		assert.Error(t, err)
	})
	t.Run("success", func(t *testing.T) {
		id := uuid.NewV4().String()
		name := "Type 001"
		category, _ := domain.NewCategory(id, name, domain.CategoryStatusInactive)

		err := category.Disable()
		assert.NoError(t, err)
	})
}

func Test_DomainCategory_ParseJson(t *testing.T) {
	t.Run("should fail if unmarshel returns error", func(t *testing.T) {
		category := new(domain.Category)
		str := "{name:Sports}"
		data := []byte(str)
		err := category.ParseJson(data)
		assert.Error(t, err)
	})
	t.Run("should fail if validation returns error", func(t *testing.T) {
		category := new(domain.Category)
		str := `{"name": "Sports"}`
		data := []byte(str)
		err := category.ParseJson(data)
		assert.Error(t, err)
	})

	t.Run("should succeed", func(t *testing.T) {
		category := new(domain.Category)
		id := uuid.NewV4().String()
		status := domain.CategoryStatusPending
		str := fmt.Sprintf(`{"id":"%s","name": "Sports", "status":"%s"}`, id, status)
		data := []byte(str)
		err := category.ParseJson(data)
		assert.NoError(t, err)
	})
}
