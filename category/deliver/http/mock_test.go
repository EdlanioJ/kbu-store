package http_test

import (
	"github.com/EdlanioJ/kbu-store/domain"
)

func testMock() *domain.Category {
	category, _ := domain.NewCategory("Store type 001")

	return category
}
