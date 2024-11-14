package product

import "github.com/lib/pq"

type ProductBase struct {
	Name        string         `json:"name" validate:"required"`
	Description string         `json:"description" validate:"required"`
	Images      pq.StringArray `json:"images"`
}

type ProductCreateRequest struct {
	*ProductBase
}

type ProductUpdateRequest struct {
	*ProductBase
}
