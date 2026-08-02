package usecase

import (
	"errors"
	"testing"

	"github.com/niviopl/clean-architecture/internal/entity"
	"github.com/stretchr/testify/suite"
)

type ListOrdersUseCaseSuite struct {
	suite.Suite
}

func TestListOrdersUseCaseSuite(t *testing.T) {
	suite.Run(t, new(ListOrdersUseCaseSuite))
}

func (s *ListOrdersUseCaseSuite) TestExecute() {
	repo := &OrderRepositoryMock{}
	orders := []entity.Order{
		{ID: "1", Price: 100, Tax: 5, FinalPrice: 105},
		{ID: "2", Price: 200, Tax: 10, FinalPrice: 210},
	}
	repo.On("FindAll").Return(orders, nil)

	usecase := NewListOrdersUseCase(repo)
	output, err := usecase.Execute()

	s.Nil(err)
	s.Len(output, 2)
	s.Equal(orders[0].ID, output[0].ID)
	s.Equal(orders[1].FinalPrice, output[1].FinalPrice)
	repo.AssertExpectations(s.T())
}

func (s *ListOrdersUseCaseSuite) TestExecute_Empty() {
	repo := &OrderRepositoryMock{}
	repo.On("FindAll").Return([]entity.Order{}, nil)

	usecase := NewListOrdersUseCase(repo)
	output, err := usecase.Execute()

	s.Nil(err)
	s.Len(output, 0)
}

func (s *ListOrdersUseCaseSuite) TestExecute_Error() {
	repo := &OrderRepositoryMock{}
	repo.On("FindAll").Return([]entity.Order{}, errors.New("db error"))

	usecase := NewListOrdersUseCase(repo)
	_, err := usecase.Execute()

	s.NotNil(err)
	s.Equal("db error", err.Error())
}
