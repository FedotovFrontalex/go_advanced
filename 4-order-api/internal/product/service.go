package product

import (
	"net/http"
	apierrors "orderApi/pkg/apiErrors"
	"orderApi/pkg/logger"
	"strconv"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ProductRepositoryInterface interface {
	Create(*Product) (*Product, error)
	GetProductById(uint) (*Product, error)
	Update(*Product) (*Product, error)
	Delete(uint) error
}

type ProductService struct {
	ProductRepository ProductRepositoryInterface
}

func NewProductService(repository ProductRepositoryInterface) *ProductService {
	return &ProductService{
		ProductRepository: repository,
	}
}

func (service *ProductService) CreateProduct(name string, description string, images pq.StringArray) (*Product, *apierrors.Error) {
	product := NewProduct(name, description, images)
	result, err := service.ProductRepository.Create(product)

	if err != nil {
		logger.Error(err)
		return nil, apierrors.NewError(http.StatusInternalServerError, err.Error())
	}

	return result, nil
}

func (service *ProductService) UpdateProduct(id uint, name string, description string, images pq.StringArray) (*Product, *apierrors.Error) {
	_, err := service.ProductRepository.GetProductById(id)

	if err != nil {
		return nil, apierrors.NewError(http.StatusNotFound, ErrNotFound)
	}

	result, err := service.ProductRepository.Update(&Product{
		Model:       gorm.Model{ID: uint(id)},
		Name:        name,
		Description: description,
		Images:      images,
	})

	if err != nil {
		return nil, apierrors.NewError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return result, nil
}

func (service *ProductService) DeleteProduct(id uint) *apierrors.Error {
	_, err := service.ProductRepository.GetProductById(id)

	if err != nil {
		return apierrors.NewError(http.StatusNotFound, ErrNotFound)
	}

	err = service.ProductRepository.Delete(uint(id))

	if err != nil {
		return apierrors.NewError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return nil
}

func (service *ProductService) GetProductById(id uint) (*Product, *apierrors.Error) {
	product, err := service.ProductRepository.GetProductById(id)

	if err != nil {
		return nil, apierrors.NewError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	return product, nil
}

func (service *ProductService) GetProductsSliceById(sliceId []string) []*Product {
	var products []*Product

	c := make(chan *Product, len(sliceId))

	for i := 0; i < len(sliceId); i++ {
		go func() {
			productId, err := strconv.ParseInt(sliceId[i], 10, 64)

			if err != nil {
				c <- nil
				return
			}

			product, err := service.ProductRepository.GetProductById(uint(productId))

			if err != nil {
				c <- nil
				return
			}

			c <- product
		}()
	}

	for i := 0; i < len(sliceId); i++ {
		product1 := <-c

		if product1 != nil {
			products = append(products, product1)
		}
	}

	return products
}
