package pg_test

import "github.com/EdlanioJ/kbu-store/domain"

func getCategory() *domain.Category {
	category, _ := domain.NewCategory("store type 001")
	return category
}
