package database

import (
	"database/sql"
	"testing"

	"github.com/niviopl/clean-architecture/internal/entity"
	"github.com/stretchr/testify/suite"

	_ "modernc.org/sqlite"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	Db *sql.DB
}

func (s *OrderRepositoryTestSuite) SetupTest() {
	db, err := sql.Open("sqlite", ":memory:")
	s.Nil(err)
	s.Db = db

	_, err = db.Exec("CREATE TABLE orders (id varchar(255), price float, tax float, final_price float)")
	s.Nil(err)
}

func (s *OrderRepositoryTestSuite) TearDownTest() {
	s.Db.Close()
}

func TestOrderRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

func (s *OrderRepositoryTestSuite) TestSave() {
	order, err := entity.NewOrder("1", 100, 5)
	s.Nil(err)

	repo := NewOrderRepository(s.Db)
	err = repo.Save(order)
	s.Nil(err)

	var orderResult entity.Order
	err = s.Db.QueryRow("SELECT id, price, tax, final_price FROM orders WHERE id = ?", order.ID).
		Scan(&orderResult.ID, &orderResult.Price, &orderResult.Tax, &orderResult.FinalPrice)
	s.Nil(err)
	s.Equal(order.ID, orderResult.ID)
	s.Equal(order.Price, orderResult.Price)
	s.Equal(order.Tax, orderResult.Tax)
	s.Equal(order.FinalPrice, orderResult.FinalPrice)
}

func (s *OrderRepositoryTestSuite) TestFindAll() {
	order1, err := entity.NewOrder("1", 100, 5)
	s.Nil(err)
	order2, err := entity.NewOrder("2", 200, 10)
	s.Nil(err)

	repo := NewOrderRepository(s.Db)
	s.Nil(repo.Save(order1))
	s.Nil(repo.Save(order2))

	orders, err := repo.FindAll()
	s.Nil(err)
	s.Len(orders, 2)
	s.Equal(order1.ID, orders[0].ID)
	s.Equal(order2.ID, orders[1].ID)
}

func (s *OrderRepositoryTestSuite) TestFindAll_Empty() {
	repo := NewOrderRepository(s.Db)
	orders, err := repo.FindAll()
	s.Nil(err)
	s.Len(orders, 0)
}
