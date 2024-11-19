package order

import (
	"orderApi/internal/product"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	Products []*product.Product `json:"products" gorm:"many2many:order_products;"`
	UserId   uint               `json:"user_id" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
