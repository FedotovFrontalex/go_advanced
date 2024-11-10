package product

import (
	"orderApi/pkg/db"
)

type ProductRepository struct {
	Database *db.Db
}

func NewProductRepository(database *db.Db) *ProductRepository {
	return &ProductRepository{
		Database: database,
	}
}

func (repo *ProductRepository) GetProductById(id uint) (*Product, error) {
	var product Product
	result := repo.Database.DB.First(&product, "id=?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &product, nil
}

func (repo *ProductRepository) Create(product *Product) (*Product, error) {
	result := repo.Database.DB.Create(product)

	if result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

func (repo *ProductRepository) Update(product *Product) (*Product, error) {
	result := repo.Database.DB.Updates(product)

	if result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

func  (repo  *ProductRepository) Delete(id uint) error {
		result := repo.Database.DB.Delete(&Product{}, "id=?", id)

		if result.Error != nil {
				return result.Error
		}

		return nil
}
