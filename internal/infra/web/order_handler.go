package web

import (
	"encoding/json"
	"net/http"

	"github.com/lmtani/learning-clean-architecture/internal/entity"
	"github.com/lmtani/learning-clean-architecture/internal/usecase"
	"github.com/lmtani/learning-clean-architecture/pkg/events"
)

type OrderHandler struct {
	EventDispatcher    events.EventDispatcherInterface
	OrderRepository    entity.OrderRepositoryInterface
	OrderCreatedEvent  events.EventInterface
	CreateOrderUseCase *usecase.CreateOrderUseCase
	ListOrdersUseCase  *usecase.ListOrdersUseCase
}

func NewOrderHandler(
	EventDispatcher events.EventDispatcherInterface,
	OrderRepository entity.OrderRepositoryInterface,
	OrderCreatedEvent events.EventInterface,
	CreateOrderUseCase *usecase.CreateOrderUseCase,
	ListOrdersUseCase *usecase.ListOrdersUseCase,
) *OrderHandler {
	return &OrderHandler{
		EventDispatcher:    EventDispatcher,
		OrderRepository:    OrderRepository,
		OrderCreatedEvent:  OrderCreatedEvent,
		CreateOrderUseCase: CreateOrderUseCase,
		ListOrdersUseCase:  ListOrdersUseCase,
	}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto usecase.OrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err := h.CreateOrderUseCase.Execute(dto)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	orders, err := h.ListOrdersUseCase.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
