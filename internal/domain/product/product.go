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
	if req.Description == "" {
		return errors.New("description is required")
	}
	if req.Type == "" {
		return errors.New("type is required")
	}
	if req.Type != "Mármore" && req.Type != "Granito" && req.Type != "Serviço" {
		return errors.New("type must be Mármore, Granito or Serviço")
	}
	if req.Price <= 0 {
		return errors.New("price must be greater than zero")
	}
	if req.Unit == "" {
		return errors.New("unit is required")
	}
	return nil
} 