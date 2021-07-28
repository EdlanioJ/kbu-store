package domain_test

import (
	"fmt"
	"testing"

	"github.com/EdlanioJ/kbu-store/app/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestCategory(t *testing.T) {
	t.Parallel()
	t.Run("new_category", func(t *testing.T) {
		category := domain.NewCategory()
		assert.NotNil(t, category)
		assert.Equal(t, category.Status, domain.CategoryStatusPending)
	})
	t.Run("parser_json", func(t *testing.T) {
		t.Parallel()
		t.Run("failure_unmarchal_returns_error", func(t *testing.T) {
			category := new(domain.Category)
			str := "{name:Sports}"
			data := []byte(str)
			err := category.ParseJson(data)
			assert.Error(t, err)
		})
		t.Run("success", func(t *testing.T) {
			category := new(domain.Category)
			id := uuid.NewV4().String()
			status := domain.CategoryStatusPending
			str := fmt.Sprintf(`{"id":"%s","name": "Sports", "status":"%s"}`, id, status)
			data := []byte(str)
			err := category.ParseJson(data)
			assert.NoError(t, err)
		})

	})
}
