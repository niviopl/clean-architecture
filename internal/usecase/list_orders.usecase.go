package usecase

type ListOrdersUseCase struct {
	OrderRepository OrderRepositoryInterface
}

func NewListOrdersUseCase(orderRepository OrderRepositoryInterface) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: orderRepository,
	}
}

func (l *ListOrdersUseCase) Execute() ([]OrderOutputDTO, error) {
	orders, err := l.OrderRepository.FindAll()
	if err != nil {
		return nil, err
	}

	output := make([]OrderOutputDTO, len(orders))
	for i, order := range orders {
		output[i] = OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		}
	}

	return output, nil
}
