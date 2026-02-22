package quote

import (
	"context"
	"time"

	quoteDomain "erp-api/internal/domain/quote"
	"erp-api/internal/utils/dbtypes"
)

type UseCaseInterface interface {
	Create(ctx context.Context, req *quoteDomain.CreateQuoteDTO) (*quoteDomain.Quote, error)
	GetByID(ctx context.Context, tenantID, id string) (*quoteDomain.Quote, error)
	Update(ctx context.Context, tenantID, id string, req *quoteDomain.UpdateQuoteDTO) (*quoteDomain.Quote, error)
	Delete(ctx context.Context, tenantID, id string) error
	List(ctx context.Context, tenantID string, limit, offset int) ([]*quoteDomain.Quote, error)
	Count(ctx context.Context, tenantID string) (int, error)
	UpdateStatus(ctx context.Context, tenantID, id string, req *quoteDomain.UpdateQuoteStatusDTO) error
}

type UseCase struct {
	quoteRepo quoteDomain.Repository
	itemRepo  quoteDomain.ItemRepository
}

func NewUseCase(quoteRepo quoteDomain.Repository, itemRepo quoteDomain.ItemRepository) UseCaseInterface {
	return &UseCase{
		quoteRepo: quoteRepo,
		itemRepo:  itemRepo,
	}
}

func (u *UseCase) Create(ctx context.Context, req *quoteDomain.CreateQuoteDTO) (*quoteDomain.Quote, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Calcular valor total
	totalValue := 0.0
	for _, item := range req.Items {
		totalValue += float64(item.Quantity) * item.Price
	}

	// Criar orçamento
	newQuote := &quoteDomain.Quote{
		TenantID:   dbtypes.UUID(req.TenantID),
		ClientID:   dbtypes.UUID(req.ClientID),
		UserID:     dbtypes.UUID(req.UserID),
		TotalValue: totalValue,
		Discount:   req.Discount,
		Status:     quoteDomain.QuoteStatusPending,
		Notes:      req.Notes,
	}

	// Se status foi fornecido, usar ele
	if req.Status != "" {
		newQuote.Status = req.Status
	}

	err := u.quoteRepo.Create(ctx, newQuote)
	if err != nil {
		return nil, err
	}

	// Criar itens do orçamento
	for _, itemDTO := range req.Items {
		item := &quoteDomain.QuoteItem{
			TenantID:  dbtypes.UUID(req.TenantID),
			QuoteID:   newQuote.ID,
			ProductID: dbtypes.UUID(itemDTO.ProductID),
			Quantity:  itemDTO.Quantity,
			UnitPrice: itemDTO.Price,
		}

		err = u.itemRepo.Create(ctx, item)
		if err != nil {
			return nil, err
		}
	}

	return newQuote, nil
}

func (u *UseCase) GetByID(ctx context.Context, tenantID, id string) (*quoteDomain.Quote, error) {
	return u.quoteRepo.GetByID(ctx, tenantID, id)
}

func (u *UseCase) Update(ctx context.Context, tenantID, id string, req *quoteDomain.UpdateQuoteDTO) (*quoteDomain.Quote, error) {
	// Buscar orçamento existente (já filtra por tenant_id)
	quote, err := u.quoteRepo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Atualizar campos
	if req.ClientID != "" {
		quote.ClientID = dbtypes.UUID(req.ClientID)
	}
	if req.UserID != "" {
		quote.UserID = dbtypes.UUID(req.UserID)
	}
	if req.Discount != nil {
		quote.Discount = *req.Discount
	}
	if req.Status != "" {
		quote.Status = req.Status
	}
	if req.Notes != "" {
		quote.Notes = req.Notes
	}

	quote.UpdatedAt = time.Now()

	err = u.quoteRepo.Update(ctx, quote)
	if err != nil {
		return nil, err
	}

	return quote, nil
}

func (u *UseCase) Delete(ctx context.Context, tenantID, id string) error {
	return u.quoteRepo.Delete(ctx, tenantID, id)
}

func (u *UseCase) List(ctx context.Context, tenantID string, limit, offset int) ([]*quoteDomain.Quote, error) {
	return u.quoteRepo.List(ctx, tenantID, limit, offset)
}

func (u *UseCase) Count(ctx context.Context, tenantID string) (int, error) {
	return u.quoteRepo.Count(ctx, tenantID)
}

func (u *UseCase) UpdateStatus(ctx context.Context, tenantID, id string, req *quoteDomain.UpdateQuoteStatusDTO) error {
	if err := req.Validate(); err != nil {
		return err
	}

	return u.quoteRepo.UpdateStatus(ctx, tenantID, id, req.Status)
}
