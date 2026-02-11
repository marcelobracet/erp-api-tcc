package product

import (
	"net/http"
	"strconv"

	productDomain "erp-api/internal/domain/product"
	productUseCase "erp-api/internal/usecase/product"
	"erp-api/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	productUseCase productUseCase.UseCaseInterface
}

func NewHandler(productUseCase productUseCase.UseCaseInterface) *Handler {
	return &Handler{
		productUseCase: productUseCase,
	}
}

func (h *Handler) Create(c *gin.Context) {
	log.Info().Msg("Create product started")

	tenantID, exists := middleware.GetTenantIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var req productDomain.CreateProductDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
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

	product, err := h.productUseCase.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Info().Msg("Create product ended")
	c.JSON(http.StatusCreated, product)
}

func (h *Handler) GetByID(c *gin.Context) {
	log.Info().Msg("Get product by ID started")

	tenantID, exists := middleware.GetTenantIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	id := c.Param("id")

	product, err := h.productUseCase.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		switch err {
		case productDomain.ErrProductNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	log.Info().Msg("Get product by ID ended")
	c.JSON(http.StatusOK, product)
}

func (h *Handler) Update(c *gin.Context) {
	log.Info().Msg("Update product started")

	tenantID, exists := middleware.GetTenantIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	id := c.Param("id")
	var req productDomain.UpdateProductDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	product, err := h.productUseCase.Update(c.Request.Context(), tenantID, id, &req)
	if err != nil {
		switch err {
		case productDomain.ErrProductNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	log.Info().Msg("Update product ended")
	c.JSON(http.StatusOK, product)
}

func (h *Handler) Delete(c *gin.Context) {
	log.Info().Msg("Delete product started")

	tenantID, exists := middleware.GetTenantIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	id := c.Param("id")

	err := h.productUseCase.Delete(c.Request.Context(), tenantID, id)
	if err != nil {
		switch err {
		case productDomain.ErrProductNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	log.Info().Msg("Delete product ended")
	c.Status(http.StatusNoContent)
}

func (h *Handler) List(c *gin.Context) {
	log.Info().Msg("List products started")

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

	products, err := h.productUseCase.List(c.Request.Context(), tenantID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	total, err := h.productUseCase.Count(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := productDomain.ProductListDTO{
		Products: make([]*productDomain.ProductDTO, len(products)),
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	}

	for i, product := range products {
		response.Products[i] = &productDomain.ProductDTO{
			ID:          product.ID,
			TenantID:    product.TenantID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			SKU:         product.SKU,
			Category:    product.Category,
			ImageURL:    product.ImageURL,
			IsActive:    product.IsActive,
			CreatedAt:   product.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:   product.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	log.Info().Msg("List products ended")
	c.JSON(http.StatusOK, response)
}

// Count conta produtos
func (h *Handler) Count(c *gin.Context) {
	log.Info().Msg("Count products started")

	tenantID, exists := middleware.GetTenantIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	count, err := h.productUseCase.Count(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Info().Msg("Count products ended")
	c.JSON(http.StatusOK, gin.H{
		"count": count,
	})
}
