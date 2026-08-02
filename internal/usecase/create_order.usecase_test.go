package usecase

import (
	"testing"

	"github.com/niviopl/clean-architecture/internal/entity"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type OrderRepositoryMock struct {
	mock.Mock
}

func (m *OrderRepositoryMock) Save(order *entity.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *OrderRepositoryMock) FindAll() ([]entity.Order, error) {
	args := m.Called()
	return args.Get(0).([]entity.Order), args.Error(1)
}

type CreateOrderUseCaseSuite struct {
	suite.Suite
}

func TestCreateOrderUseCaseSuite(t *testing.T) {
	suite.Run(t, new(CreateOrderUseCaseSuite))
}

func (s *CreateOrderUseCaseSuite) TestExecute() {
	repo := &OrderRepositoryMock{}
	repo.On("Save", mock.Anything).Return(nil)

	usecase := NewCreateOrderUseCase(repo)
	input := OrderInputDTO{ID: "1", Price: 100, Tax: 5}

	output, err := usecase.Execute(input)

	s.Nil(err)
	s.Equal(input.ID, output.ID)
	s.Equal(input.Price, output.Price)
	s.Equal(input.Tax, output.Tax)
	s.Equal(105.0, output.FinalPrice)
	repo.AssertExpectations(s.T())
	repo.AssertNumberOfCalls(s.T(), "Save", 1)
}

func (s *CreateOrderUseCaseSuite) TestExecute_InvalidInput() {
	repo := &OrderRepositoryMock{}

	usecase := NewCreateOrderUseCase(repo)
	input := OrderInputDTO{ID: "", Price: 100, Tax: 5}

	_, err := usecase.Execute(input)

	s.Equal(entity.ErrIDIsRequired, err)
	repo.AssertNotCalled(s.T(), "Save")
}
