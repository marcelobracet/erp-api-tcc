package product

import "errors"

var (
	ErrProductNotFound     = errors.New("product not found")
	ErrProductAlreadyExists = errors.New("product already exists")
	ErrInvalidProductType  = errors.New("invalid product type")
)

func (req *CreateProductDTO) Validate() error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if req.Price <= 0 {
		return errors.New("price must be greater than zero")
	}
	return nil
} 