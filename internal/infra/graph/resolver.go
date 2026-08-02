package graph

// THIS CODE WILL BE UPDATED WITH SCHEMA CHANGES. PREVIOUS IMPLEMENTATION FOR SCHEMA CHANGES WILL BE KEPT IN THE COMMENT SECTION. IMPLEMENTATION FOR UNCHANGED SCHEMA WILL BE KEPT.

import (
	"context"

	"github.com/niviopl/clean-architecture/internal/infra/graph/generated"
	"github.com/niviopl/clean-architecture/internal/infra/graph/model"
	"github.com/niviopl/clean-architecture/internal/usecase"
)

type Resolver struct {
	CreateOrderUseCase *usecase.CreateOrderUseCase
	ListOrdersUseCase  *usecase.ListOrdersUseCase
}

// CreateOrder is the resolver for the createOrder field.
func (r *mutationResolver) CreateOrder(ctx context.Context, input model.OrderInput) (*model.Order, error) {
	output, err := r.CreateOrderUseCase.Execute(usecase.OrderInputDTO{
		ID:    input.ID,
		Price: input.Price,
		Tax:   input.Tax,
	})
	if err != nil {
		return nil, err
	}

	return &model.Order{
		ID:         output.ID,
		Price:      output.Price,
		Tax:        output.Tax,
		FinalPrice: output.FinalPrice,
	}, nil
}

// ListOrders is the resolver for the ListOrders field.
func (r *queryResolver) ListOrders(ctx context.Context) ([]*model.Order, error) {
	output, err := r.ListOrdersUseCase.Execute()
	if err != nil {
		return nil, err
	}

	orders := make([]*model.Order, len(output))
	for i, order := range output {
		orders[i] = &model.Order{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		}
	}

	return orders, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type (
	mutationResolver struct{ *Resolver }
	queryResolver    struct{ *Resolver }
)
