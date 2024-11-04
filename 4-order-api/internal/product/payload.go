package product

import "github.com/lib/pq"

type AddProductRequest struct {
		Name string `json:"name" validate:"required"`
		Description string `json:"description" validate:"required"`
		Images pq.StringArray `json:"images"` 
}
