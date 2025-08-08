package client

import (
	"context"
	"time"

	clientDomain "erp-api/internal/domain/client"

	"github.com/google/uuid"
)

type UseCaseInterface interface {
	Create(ctx context.Context, req *clientDomain.CreateClientDTO) (*clientDomain.Client, error)
	GetByID(ctx context.Context, id string) (*clientDomain.Client, error)
	Update(ctx context.Context, id string, req *clientDomain.UpdateClientDTO) (*clientDomain.Client, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*clientDomain.Client, error)
	Count(ctx context.Context) (int, error)
}

type UseCase struct {
	clientRepo clientDomain.Repository
}

func NewUseCase(clientRepo clientDomain.Repository) UseCaseInterface {
	return &UseCase{
		clientRepo: clientRepo,
	}
}

func (u *UseCase) Create(ctx context.Context, req *clientDomain.CreateClientDTO) (*clientDomain.Client, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Verificar se cliente já existe
	existingClient, err := u.clientRepo.GetByDocument(ctx, req.Document)
	if err != nil && err != clientDomain.ErrClientNotFound {
		return nil, err
	}

	if existingClient != nil {
		return nil, clientDomain.ErrClientAlreadyExists
	}

	// Criar cliente
	newClient := &clientDomain.Client{
		ID:           uuid.New().String(),
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		Document:     req.Document,
		DocumentType: req.DocumentType,
		Address:      req.Address,
		City:         req.City,
		State:        req.State,
		ZipCode:      req.ZipCode,
		IsActive:     true,
	}

	err = u.clientRepo.Create(ctx, newClient)
	if err != nil {
		return nil, err
	}

	return newClient, nil
}

func (u *UseCase) GetByID(ctx context.Context, id string) (*clientDomain.Client, error) {
	return u.clientRepo.GetByID(ctx, id)
}

func (u *UseCase) Update(ctx context.Context, id string, req *clientDomain.UpdateClientDTO) (*clientDomain.Client, error) {
	// Buscar cliente existente
	client, err := u.clientRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Atualizar campos
	if req.Name != "" {
		client.Name = req.Name
	}
	if req.Email != "" {
		client.Email = req.Email
	}
	if req.Phone != "" {
		client.Phone = req.Phone
	}
	if req.Document != "" {
		client.Document = req.Document
	}
	if req.DocumentType != "" {
		client.DocumentType = req.DocumentType
	}
	if req.Address != "" {
		client.Address = req.Address
	}
	if req.City != "" {
		client.City = req.City
	}
	if req.State != "" {
		client.State = req.State
	}
	if req.ZipCode != "" {
		client.ZipCode = req.ZipCode
	}
	if req.IsActive != nil {
		client.IsActive = *req.IsActive
	}

	client.UpdatedAt = time.Now()

	err = u.clientRepo.Update(ctx, client)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (u *UseCase) Delete(ctx context.Context, id string) error {
	return u.clientRepo.Delete(ctx, id)
}

func (u *UseCase) List(ctx context.Context, limit, offset int) ([]*clientDomain.Client, error) {
	return u.clientRepo.List(ctx, limit, offset)
}

func (u *UseCase) Count(ctx context.Context) (int, error) {
	return u.clientRepo.Count(ctx)
} 