package usecase

import "github.com/niviopl/clean-architecture/internal/entity"

type CreateOrderUseCase struct {
	OrderRepository OrderRepositoryInterface
}

func NewCreateOrderUseCase(orderRepository OrderRepositoryInterface) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository: orderRepository,
	}
}

func (c *CreateOrderUseCase) Execute(input OrderInputDTO) (OrderOutputDTO, error) {
	order, err := entity.NewOrder(input.ID, input.Price, input.Tax)
	if err != nil {
		return OrderOutputDTO{}, err
	}

	if err := c.OrderRepository.Save(order); err != nil {
		return OrderOutputDTO{}, err
	}

	return OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}, nil
}
