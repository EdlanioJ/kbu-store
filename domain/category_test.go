package domain_test

import (
	"testing"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewCategory(t *testing.T) {
	t.Run("should fail on validation", func(t *testing.T) {
		is := assert.New(t)

		category, err := domain.NewCategory("")

		is.Nil(category)
		is.Error(err)
	})

	t.Run("should succeed", func(t *testing.T) {
		is := assert.New(t)

		name := "Type 001"
		category, err := domain.NewCategory(name)

		is.NotNil(category)
		is.Nil(err)
		is.Equal(category.Name, name)
	})
}

func Test_DomainCategory_Activate(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		name := "category001"
		category, _ := domain.NewCategory(name)
		category.ID = ""

		err := category.Activate()
		assert.Error(t, err)
	})
	t.Run("success", func(t *testing.T) {
		name := "category001"
		category, _ := domain.NewCategory(name)

		err := category.Activate()
		assert.NoError(t, err)
	})
}

func Test_DomainCategory_Disable(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		name := "category001"
		category, _ := domain.NewCategory(name)
		category.ID = ""

		err := category.Disable()
		assert.Error(t, err)
	})
	t.Run("success", func(t *testing.T) {
		name := "category001"
		category, _ := domain.NewCategory(name)

		err := category.Disable()
		assert.NoError(t, err)
	})
}
