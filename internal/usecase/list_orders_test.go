package usecase

import (
	"testing"

	"github.com/lmtani/learning-clean-architecture/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOrderRepository - Não sei qual o melhor pacote para colocar o mock. Por hora deixei aqui
// mas penso que talvez pudesse ser no 'database'.
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) ListAll() ([]*entity.Order, error) {
	args := m.Called()
	return args.Get(0).([]*entity.Order), args.Error(1)
}

func (m *MockOrderRepository) Save(*entity.Order) error {
	args := m.Called()
	return args.Error(0)
}

func TestListOrdersUseCase_Execute(t *testing.T) {
	tests := []struct {
		name           string
		mockReturn     []*entity.Order
		mockError      error
		expectedOutput []*OrderOutputDTO
		expectedError  error
	}{
		{
			name: "orders exist",
			mockReturn: []*entity.Order{
				{
					ID:         "1",
					Price:      10.0,
					Tax:        2.0,
					FinalPrice: 12.0,
				},
				{
					ID:         "2",
					Price:      20.0,
					Tax:        5.0,
					FinalPrice: 25.0,
				},
			},
			mockError: nil,
			expectedOutput: []*OrderOutputDTO{
				{
					ID:         "1",
					Price:      10.0,
					Tax:        2.0,
					FinalPrice: 12.0,
				},
				{
					ID:         "2",
					Price:      20.0,
					Tax:        5.0,
					FinalPrice: 25.0,
				},
			},
			expectedError: nil,
		},
		{
			name:           "no orders",
			mockReturn:     []*entity.Order{},
			mockError:      nil,
			expectedOutput: []*OrderOutputDTO{},
			expectedError:  nil,
		},
		{
			name:           "error",
			mockReturn:     nil,
			mockError:      assert.AnError,
			expectedOutput: nil,
			expectedError:  assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Arrange
			mockRepo := &MockOrderRepository{}
			listOrdersUseCase := NewListOrdersUseCase(mockRepo)

			mockRepo.On("ListAll").Return(tt.mockReturn, tt.mockError).Once()

			// Act
			outputDTOs, err := listOrdersUseCase.Execute()

			// Assert
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, len(tt.expectedOutput), len(outputDTOs))

			for i, expectedOrder := range tt.expectedOutput {
				assert.Equal(t, expectedOrder.ID, outputDTOs[i].ID)
				assert.Equal(t, expectedOrder.Price, outputDTOs[i].Price)
				assert.Equal(t, expectedOrder.Tax, outputDTOs[i].Tax)
				assert.Equal(t, expectedOrder.FinalPrice, outputDTOs[i].FinalPrice)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
