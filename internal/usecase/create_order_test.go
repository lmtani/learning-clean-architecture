package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/lmtani/learning-clean-architecture/pkg/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Não sei qual o melhor pacote para colocar o mock.

// MockEvent - Sei que existe o TestEvent no pacote events, mas não sei se seria o caso de usar ele
// ou criar um novo mock
type MockEvent struct {
	mock.Mock
}

// GetDateTime - existe na interface utilizada mas não é utilizada no caso de uso
func (me *MockEvent) GetDateTime() time.Time {
	panic("unimplemented")
}

func (me *MockEvent) GetName() string {
	args := me.Called()
	return args.String(0)
}

func (me *MockEvent) GetPayload() interface{} {
	args := me.Called()
	return args.Get(0)
}

func (me *MockEvent) SetPayload(payload interface{}) {
	me.Called(payload)
}

type MockEventDispatcher struct {
	mock.Mock
}

// Has - existe na interface utilizada mas não é utilizada no caso de uso
func (md *MockEventDispatcher) Has(eventName string, handler events.EventHandlerInterface) bool {
	panic("unimplemented")
}

// Register
func (md *MockEventDispatcher) Register(eventName string, handler events.EventHandlerInterface) error {
	panic("unimplemented")
}

// Remove
func (md *MockEventDispatcher) Remove(eventName string, handler events.EventHandlerInterface) error {
	panic("unimplemented")
}

func (md *MockEventDispatcher) Dispatch(event events.EventInterface) error {
	args := md.Called(event)
	return args.Error(0)
}

func (md *MockEventDispatcher) Clear() {
	md.Called()
}

func TestCreateOrderUseCase_Execute(t *testing.T) {
	tests := []struct {
		name               string
		input              OrderInputDTO
		expectedOutput     OrderOutputDTO
		calculatePriceErr  bool
		mockSaveError      error
		expectEventPayload bool
		expectDispatch     bool
		expectedError      error
	}{
		{
			name: "success",
			input: OrderInputDTO{
				ID:    "1",
				Price: 100,
				Tax:   10,
			},
			expectedOutput: OrderOutputDTO{
				ID:         "1",
				Price:      100,
				Tax:        10,
				FinalPrice: 110, // Price + Tax
			},
			calculatePriceErr:  false,
			mockSaveError:      nil,
			expectEventPayload: true,
			expectDispatch:     true,
			expectedError:      nil,
		},
		{
			name: "error saving order",
			input: OrderInputDTO{
				ID:    "3",
				Price: 100,
				Tax:   10,
			},
			expectedOutput:     OrderOutputDTO{},
			calculatePriceErr:  false,
			mockSaveError:      errors.New("repository error"),
			expectEventPayload: false,
			expectDispatch:     false,
			expectedError:      errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockOrderRepository{}
			mockEvent := &MockEvent{}
			mockDispatcher := &MockEventDispatcher{}
			mockRepo.On("Save", mock.Anything).Return(tt.mockSaveError).Once()

			if tt.expectEventPayload {
				mockEvent.On("SetPayload", tt.expectedOutput).Once()
			}
			if tt.expectDispatch {
				mockDispatcher.On("Dispatch", mockEvent).Return(tt.expectedError).Once()
			}

			createOrderUseCase := NewCreateOrderUseCase(mockRepo, mockEvent, mockDispatcher)

			// Act
			output, err := createOrderUseCase.Execute(tt.input)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Equal(t, tt.expectedOutput, output)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, output)
			}

			mockRepo.AssertExpectations(t)
			mockEvent.AssertExpectations(t)
			mockDispatcher.AssertExpectations(t)
		})
	}
}
