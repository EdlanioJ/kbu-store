package usecase

import (
	"context"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type fillCategoryDetails struct {
	getCategoryRepo domain.GetCategoryByIDRepository
}

func newFillCategoryDetails(getRepo domain.GetCategoryByIDRepository) *fillCategoryDetails {
	return &fillCategoryDetails{
		getCategoryRepo: getRepo,
	}
}

func (u *fillCategoryDetails) exec(c context.Context, data []*domain.Store) ([]*domain.Store, error) {
	g, ctx := errgroup.WithContext(c)

	mapCategories := map[string]*domain.Category{}
	for _, store := range data {
		mapCategories[store.Category.ID] = &domain.Category{}
	}

	chanCategory := make(chan *domain.Category)
	for categoryID := range mapCategories {
		categoryID := categoryID
		g.Go(func() error {
			res, err := u.getCategoryRepo.Exec(ctx, categoryID)
			if err != nil {
				return err
			}
			chanCategory <- res
			return nil
		})
	}

	go func() {
		err := g.Wait()
		if err != nil {
			logrus.Error(err)
		}

		close(chanCategory)
	}()

	for category := range chanCategory {
		if category != (&domain.Category{}) {
			mapCategories[category.ID] = category
		}
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	for index, item := range data {
		if a, ok := mapCategories[item.Category.ID]; ok {
			data[index].Category = a
		}
	}
	return data, nil
}
