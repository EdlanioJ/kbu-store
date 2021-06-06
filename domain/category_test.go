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
