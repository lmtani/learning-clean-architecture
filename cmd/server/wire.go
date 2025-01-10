//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"

	"github.com/google/wire"
	"github.com/lmtani/learning-clean-architecture/internal/entity"
	"github.com/lmtani/learning-clean-architecture/internal/infra/database"
	"github.com/lmtani/learning-clean-architecture/internal/infra/web"
	"github.com/lmtani/learning-clean-architecture/internal/usecase"
)

var setOrderRepositoryDependency = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
)

// var setEventDispatcherDependency = wire.NewSet(
// 	events.NewEventDispatcher,
// 	event.NewOrderCreated,
// 	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
// 	wire.Bind(new(events.EventDispatcherInterface), new(*events.EventDispatcher)),
// )

// var setOrderCreatedEvent = wire.NewSet(
// 	event.NewOrderCreated,
// 	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
// )

func NewCreateOrderUseCase(db *sql.DB) *usecase.CreateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		usecase.NewCreateOrderUseCase,
	)
	return &usecase.CreateOrderUseCase{}
}

func NewWebOrderHandler(db *sql.DB) *web.OrderHandler {
	wire.Build(
		setOrderRepositoryDependency,
		web.NewWebOrderHandler,
	)
	return &web.OrderHandler{}
}
