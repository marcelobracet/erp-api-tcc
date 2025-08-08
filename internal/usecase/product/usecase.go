package product

import (
	"context"
	"time"

	productDomain "erp-api/internal/domain/product"

	"github.com/google/uuid"
)

type UseCaseInterface interface {
	Create(ctx context.Context, req *productDomain.CreateProductDTO) (*productDomain.Product, error)
	GetByID(ctx context.Context, id string) (*productDomain.Product, error)
	Update(ctx context.Context, id string, req *productDomain.UpdateProductDTO) (*productDomain.Product, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*productDomain.Product, error)
	Count(ctx context.Context) (int, error)
}

type UseCase struct {
	productRepo productDomain.Repository
}

func NewUseCase(productRepo productDomain.Repository) UseCaseInterface {
	return &UseCase{
		productRepo: productRepo,
	}
}

func (u *UseCase) Create(ctx context.Context, req *productDomain.CreateProductDTO) (*productDomain.Product, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Criar produto
	newProduct := &productDomain.Product{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Price:       req.Price,
		Unit:        req.Unit,
		IsActive:    true,
	}

	err := u.productRepo.Create(ctx, newProduct)
	if err != nil {
		return nil, err
	}

	return newProduct, nil
}

func (u *UseCase) GetByID(ctx context.Context, id string) (*productDomain.Product, error) {
	return u.productRepo.GetByID(ctx, id)
}

func (u *UseCase) Update(ctx context.Context, id string, req *productDomain.UpdateProductDTO) (*productDomain.Product, error) {
	// Buscar produto existente
	product, err := u.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Atualizar campos
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Type != "" {
		product.Type = req.Type
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Unit != "" {
		product.Unit = req.Unit
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	product.UpdatedAt = time.Now()

	err = u.productRepo.Update(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (u *UseCase) Delete(ctx context.Context, id string) error {
	return u.productRepo.Delete(ctx, id)
}

func (u *UseCase) List(ctx context.Context, limit, offset int) ([]*productDomain.Product, error) {
	return u.productRepo.List(ctx, limit, offset)
}

func (u *UseCase) Count(ctx context.Context) (int, error) {
	return u.productRepo.Count(ctx)
} 