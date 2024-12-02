package product_test

import (
	"errors"
	"net/http"
	"orderApi/internal/product"
	"testing"

	"gorm.io/gorm"
)

type MockProductRepository struct{}

func (r *MockProductRepository) Create(data *product.Product) (*product.Product, error) {
	return &product.Product{}, nil
}

func (r *MockProductRepository) GetProductById(id uint) (*product.Product, error) {
	if id == EXISTED_PRODUCT_ID {
		return &product.Product{
			Model: gorm.Model{ID: EXISTED_PRODUCT_ID},
		}, nil
	}

	return nil, errors.New("some error")
}

func (r *MockProductRepository) Update(data *product.Product) (*product.Product, error) {
	if data.ID != EXISTED_PRODUCT_ID {
		return nil, errors.New("product not found")
	}

	return &product.Product{
		Model:       gorm.Model{ID: data.ID},
		Name:        data.Name,
		Description: data.Description,
		Images:      data.Images,
	}, nil
}

func (r *MockProductRepository) Delete(id uint) error {
	return nil
}

func bootstrap() *product.ProductService {
	productRepository := MockProductRepository{}

	return &product.ProductService{
		ProductRepository: &productRepository,
	}
}

func TestCreateProductSuccess(t *testing.T) {
	service := bootstrap()

	productData, _ := service.CreateProduct("product", "description", nil)

	if productData == nil {
		t.Fatal("Product data not to be nil")
	}
}

func TestUpdateProductSuccess(t *testing.T) {
	service := bootstrap()

	updatedProductName := "UpdatedName"
	updatedProductDescription := "UpdatedDescription"

	updatedProduct, err := service.UpdateProduct(EXISTED_PRODUCT_ID, updatedProductName, updatedProductDescription, nil)

	if err != nil {
		t.Fatal(err)
	}

	if updatedProduct.ID != EXISTED_PRODUCT_ID {
		t.Fatalf("Expected product id %d got %d", EXISTED_PRODUCT_ID, updatedProduct.ID)
	}

	if updatedProduct.Name != updatedProductName {
		t.Fatalf("Expected product name %s got %s", updatedProductName, updatedProduct.Name)
	}

	if updatedProduct.Description != updatedProductDescription {
		t.Fatalf("Expected product description %s got %s", updatedProductDescription, updatedProduct.Description)
	}
}

func TestUpdateProductFailedWithNoExistedProduct(t *testing.T) {
	service := bootstrap()

	updatedProductName := "UpdatedName"
	updatedProductDescription := "UpdatedDescription"

	updatedProduct, err := service.UpdateProduct(NO_EXISTED_PRODUCTID, updatedProductName, updatedProductDescription, nil)

	if err == nil {
		t.Fatal("Expected err")
		return
	}

	if err.GetStatus() != http.StatusNotFound {
		t.Fatalf("Expected err status %d got %d", http.StatusNotFound, err.GetStatus())
	}

	if err.Error() != product.ErrNotFound {
		t.Fatalf("Expected err text %s got %s", product.ErrNotFound, err.Error())
	}

	if updatedProduct != nil {
		t.Fatal("product result nust be nil")
	}
}

func TestDeleteProductSuccess(t *testing.T) {
	service := bootstrap()

	err := service.DeleteProduct(EXISTED_PRODUCT_ID)

	if err != nil {
		t.Fatal("Error must be nil")
	}
}

func TestDeleteProductFailedWithNoExistProduct(t *testing.T) {
	service := bootstrap()

	err := service.DeleteProduct(NO_EXISTED_PRODUCTID)

	if err.GetStatus() != http.StatusNotFound {
		t.Fatalf("Expected Error status %d got %d", http.StatusNotFound, err.GetStatus())
	}
}
