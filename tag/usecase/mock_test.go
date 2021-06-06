package usecase_test

import "github.com/EdlanioJ/kbu-store/domain"

func testMock() *domain.Tag {
	return &domain.Tag{
		Name:  "tag001",
		Count: 2,
	}
}
