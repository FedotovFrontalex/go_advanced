package order

import (
	"orderApi/internal/product"
	"orderApi/pkg/db"
)

type OrderRepository struct {
	Db *db.Db
}

func NewOrderRepository(db *db.Db) *OrderRepository {
	return &OrderRepository{
		Db: db,
	}
}

func (repo *OrderRepository) CreateOrder(products []*product.Product, userId uint) (*Order, error) {
	order := &Order{
		UserId:   userId,
		Products: products,
	}

	result := repo.Db.Create(order)

	if result.Error != nil {
		return nil, result.Error
	}

	return order, nil
}

func (repo *OrderRepository) GetById(id uint) (*Order, error) {
	var order Order

	result := repo.Db.
		Preload("Products").
		Find(&order, "id=?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &order, nil
}

func (repo *OrderRepository) GetByUserId(id uint) []Order {
	var order []Order

	repo.Db.Table("orders").
		Preload("Products").
		Order("updated_at DESC").
		Find(&order, "user_id=?", id)

	return order
}
