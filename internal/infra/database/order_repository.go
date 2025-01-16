package database

import (
	"context"

	"github.com/lmtani/learning-clean-architecture/internal/entity"
	"github.com/lmtani/learning-clean-architecture/internal/infra/database/psql"
)

type OrderRepository struct {
	Queries *psql.Queries
}

func NewOrderRepository(queries *psql.Queries) *OrderRepository {
	return &OrderRepository{Queries: queries}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	err := r.Queries.Save(
		context.Background(),
		psql.SaveParams{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) ListAll() ([]*entity.Order, error) {
	orders, err := r.Queries.ListAll(context.Background())
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Order, 0)
	for _, order := range orders {
		result = append(result, &entity.Order{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		})
	}
	return result, nil
}
