package usecase_test

import (
	"github.com/EdlanioJ/kbu-store/domain"
)

func testMock() *domain.Category {

	mockStorType, _ := domain.NewCategory("Store type 001")

	return mockStorType
}
