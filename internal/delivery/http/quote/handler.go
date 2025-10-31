package quote

import (
	"net/http"
	"strconv"

	quoteDomain "erp-api/internal/domain/quote"
	quoteUseCase "erp-api/internal/usecase/quote"
	"erp-api/pkg/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	quoteUseCase quoteUseCase.UseCaseInterface
}

func NewHandler(quoteUseCase quoteUseCase.UseCaseInterface) *Handler {
	return &Handler{
		quoteUseCase: quoteUseCase,
	}
}

// Create cria um novo orçamento
func (h *Handler) Create(c *gin.Context) {
	tenantID, exists := middleware.GetTenantIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var req quoteDomain.CreateQuoteDTO
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validar que o tenant_id do request corresponde ao tenant_id do usuário autenticado
	if req.TenantID != tenantID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Tenant ID mismatch",
		})
		return
	}
	
	quote, err := h.quoteUseCase.Create(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case quoteDomain.ErrInvalidItems:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Quote must have at least one item",
			})
		case quoteDomain.ErrInvalidQuoteStatus:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid quote status",
			})
		case quoteDomain.ErrInvalidDate:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid date",
			})
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		return
	}
	
	c.JSON(http.StatusCreated, quote)
}

// GetByID obtém um orçamento por ID
func (h *Handler) GetByID(c *gin.Context) {
	tenantID, exists := middleware.GetTenantIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	id := c.Param("id")
	
	quote, err := h.quoteUseCase.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		switch err {
		case quoteDomain.ErrQuoteNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Quote not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}
	
	c.JSON(http.StatusOK, quote)
}

// Update atualiza um orçamento
func (h *Handler) Update(c *gin.Context) {
	tenantID, exists := middleware.GetTenantIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	id := c.Param("id")
	var req quoteDomain.UpdateQuoteDTO
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	quote, err := h.quoteUseCase.Update(c.Request.Context(), tenantID, id, &req)
	if err != nil {
		switch err {
		case quoteDomain.ErrQuoteNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Quote not found",
			})
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		return
	}
	
	c.JSON(http.StatusOK, quote)
}

// Delete deleta um orçamento
func (h *Handler) Delete(c *gin.Context) {
	tenantID, exists := middleware.GetTenantIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	id := c.Param("id")
	
	err := h.quoteUseCase.Delete(c.Request.Context(), tenantID, id)
	if err != nil {
		switch err {
		case quoteDomain.ErrQuoteNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Quote not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}
	
	c.Status(http.StatusNoContent)
}

// List lista orçamentos
func (h *Handler) List(c *gin.Context) {
	tenantID, exists := middleware.GetTenantIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid limit parameter",
		})
		return
	}
	
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid offset parameter",
		})
		return
	}
	
	quotes, err := h.quoteUseCase.List(c.Request.Context(), tenantID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	total, err := h.quoteUseCase.Count(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	response := quoteDomain.QuoteListDTO{
		Quotes: make([]*quoteDomain.QuoteDTO, len(quotes)),
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}
	
	for i, quote := range quotes {
		response.Quotes[i] = &quoteDomain.QuoteDTO{
			ID:             quote.ID,
			TenantID:       quote.TenantID,
			ClientID:       quote.ClientID,
			UserID:         quote.UserID,
			TotalValue:     quote.TotalValue,
			Discount:       quote.Discount,
			Status:         quote.Status,
			ConversionRate: quote.ConversionRate,
			Notes:          quote.Notes,
			CreatedAt:      quote.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:      quote.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	
	c.JSON(http.StatusOK, response)
}

// Count conta orçamentos
func (h *Handler) Count(c *gin.Context) {
	tenantID, exists := middleware.GetTenantIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	count, err := h.quoteUseCase.Count(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"count": count,
	})
}

// UpdateStatus atualiza o status de um orçamento
func (h *Handler) UpdateStatus(c *gin.Context) {
	tenantID, exists := middleware.GetTenantIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	id := c.Param("id")
	var req quoteDomain.UpdateQuoteStatusDTO
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	err := h.quoteUseCase.UpdateStatus(c.Request.Context(), tenantID, id, &req)
	if err != nil {
		switch err {
		case quoteDomain.ErrQuoteNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Quote not found",
			})
		case quoteDomain.ErrInvalidQuoteStatus:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid quote status",
			})
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Quote status updated successfully",
	})
} 