package graph

import "github.com/lmtani/learning-clean-architecture/internal/usecase"

type Resolver struct {
	CreateOrderUseCase usecase.CreateOrderUseCase
}
