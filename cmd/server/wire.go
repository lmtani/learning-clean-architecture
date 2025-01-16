//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"

	"github.com/google/wire"

	"github.com/lmtani/learning-clean-architecture/internal/entity"
	"github.com/lmtani/learning-clean-architecture/internal/infra/database"
	"github.com/lmtani/learning-clean-architecture/internal/infra/event"
	"github.com/lmtani/learning-clean-architecture/internal/infra/web"
	"github.com/lmtani/learning-clean-architecture/internal/usecase"
	"github.com/lmtani/learning-clean-architecture/pkg/events"
)

var setOrderRepositoryDependency = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
)

var setOrderCreatedEvent = wire.NewSet(
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
)

func NewWebOrderHandler(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *web.OrderHandler {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		web.NewOrderHandler,
		usecase.NewCreateOrderUseCase,
		usecase.NewListOrdersUseCase,
	)
	return &web.OrderHandler{}
}

func NewCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.CreateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		usecase.NewCreateOrderUseCase,
	)
	return &usecase.CreateOrderUseCase{}
}

func NewListOrdersUseCase(db *sql.DB) *usecase.ListOrdersUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		usecase.NewListOrdersUseCase,
	)
	return &usecase.ListOrdersUseCase{}
}
