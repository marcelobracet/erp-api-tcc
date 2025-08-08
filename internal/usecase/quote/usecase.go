package quote

import (
	"context"
	"time"

	quoteDomain "erp-api/internal/domain/quote"

	"github.com/google/uuid"
)

type UseCaseInterface interface {
	Create(ctx context.Context, req *quoteDomain.CreateQuoteDTO) (*quoteDomain.Quote, error)
	GetByID(ctx context.Context, id string) (*quoteDomain.Quote, error)
	Update(ctx context.Context, id string, req *quoteDomain.UpdateQuoteDTO) (*quoteDomain.Quote, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*quoteDomain.Quote, error)
	Count(ctx context.Context) (int, error)
	UpdateStatus(ctx context.Context, id string, req *quoteDomain.UpdateQuoteStatusDTO) error
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
		// Aqui você faria uma consulta ao produto para pegar o preço
		// Por simplicidade, vamos assumir que o preço vem do frontend
		totalValue += item.Quantity * 100.0 // Preço padrão, deve vir do produto
	}

	// Criar orçamento
	newQuote := &quoteDomain.Quote{
		ID:         uuid.New().String(),
		ClientID:   req.ClientID,
		TotalValue: totalValue,
		Status:     quoteDomain.QuoteStatusPending,
		Date:       req.Date,
		ValidUntil: req.ValidUntil,
		Notes:      req.Notes,
		IsActive:   true,
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
			ID:         uuid.New().String(),
			QuoteID:    newQuote.ID,
			ProductID:  itemDTO.ProductID,
			Quantity:   itemDTO.Quantity,
			UnitPrice:  100.0, // Deve vir do produto
			TotalPrice: itemDTO.Quantity * 100.0,
		}

		err = u.itemRepo.Create(ctx, item)
		if err != nil {
			return nil, err
		}
	}

	return newQuote, nil
}

func (u *UseCase) GetByID(ctx context.Context, id string) (*quoteDomain.Quote, error) {
	return u.quoteRepo.GetByID(ctx, id)
}

func (u *UseCase) Update(ctx context.Context, id string, req *quoteDomain.UpdateQuoteDTO) (*quoteDomain.Quote, error) {
	// Buscar orçamento existente
	quote, err := u.quoteRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Atualizar campos
	if req.ClientID != "" {
		quote.ClientID = req.ClientID
	}
	if req.Status != "" {
		quote.Status = req.Status
	}
	if req.Date != nil {
		quote.Date = *req.Date
	}
	if req.ValidUntil != nil {
		quote.ValidUntil = *req.ValidUntil
	}
	if req.Notes != "" {
		quote.Notes = req.Notes
	}
	if req.IsActive != nil {
		quote.IsActive = *req.IsActive
	}

	quote.UpdatedAt = time.Now()

	err = u.quoteRepo.Update(ctx, quote)
	if err != nil {
		return nil, err
	}

	return quote, nil
}

func (u *UseCase) Delete(ctx context.Context, id string) error {
	return u.quoteRepo.Delete(ctx, id)
}

func (u *UseCase) List(ctx context.Context, limit, offset int) ([]*quoteDomain.Quote, error) {
	return u.quoteRepo.List(ctx, limit, offset)
}

func (u *UseCase) Count(ctx context.Context) (int, error) {
	return u.quoteRepo.Count(ctx)
}

func (u *UseCase) UpdateStatus(ctx context.Context, id string, req *quoteDomain.UpdateQuoteStatusDTO) error {
	if err := req.Validate(); err != nil {
		return err
	}

	return u.quoteRepo.UpdateStatus(ctx, id, req.Status)
} 