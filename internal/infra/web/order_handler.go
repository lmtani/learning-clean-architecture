package web

import (
	"encoding/json"
	"net/http"

	"github.com/lmtani/learning-clean-architecture/internal/entity"
	"github.com/lmtani/learning-clean-architecture/internal/usecase"
)

type OrderHandler struct {
	// EventDispatcher   events.EventDispatcherInterface
	OrderRepository entity.OrderRepositoryInterface
	// OrderCreatedEvent events.EventInterface
}

func NewWebOrderHandler(
	// EventDispatcher events.EventDispatcherInterface,
	OrderRepository entity.OrderRepositoryInterface,
	// OrderCreatedEvent events.EventInterface,
) *OrderHandler {
	return &OrderHandler{
		// EventDispatcher:   EventDispatcher,
		OrderRepository: OrderRepository,
		// OrderCreatedEvent: OrderCreatedEvent,
	}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto usecase.OrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createOrder := usecase.NewCreateOrderUseCase(h.OrderRepository)
	output, err := createOrder.Execute(dto)
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
