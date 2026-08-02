package usecase

import "github.com/niviopl/clean-architecture/internal/entity"

type OrderRepositoryInterface interface {
	Save(order *entity.Order) error
	FindAll() ([]entity.Order, error)
}
