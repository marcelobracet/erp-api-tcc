package product

import (
	"net/http"
	"strconv"

	productDomain "erp-api/internal/domain/product"
	productUseCase "erp-api/internal/usecase/product"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	productUseCase productUseCase.UseCaseInterface
}

func NewHandler(productUseCase productUseCase.UseCaseInterface) *Handler {
	return &Handler{
		productUseCase: productUseCase,
	}
}

// Create cria um novo produto
func (h *Handler) Create(c *gin.Context) {
	var req productDomain.CreateProductDTO
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
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
	
	c.JSON(http.StatusCreated, product)
}

// GetByID obtém um produto por ID
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	
	product, err := h.productUseCase.GetByID(c.Request.Context(), id)
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
	
	c.JSON(http.StatusOK, product)
}

// Update atualiza um produto
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var req productDomain.UpdateProductDTO
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	product, err := h.productUseCase.Update(c.Request.Context(), id, &req)
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
	
	c.JSON(http.StatusOK, product)
}

// Delete deleta um produto
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	
	err := h.productUseCase.Delete(c.Request.Context(), id)
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
	
	c.Status(http.StatusNoContent)
}

// List lista produtos
func (h *Handler) List(c *gin.Context) {
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
	
	products, err := h.productUseCase.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	total, err := h.productUseCase.Count(c.Request.Context())
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
	
	c.JSON(http.StatusOK, response)
}

// Count conta produtos
func (h *Handler) Count(c *gin.Context) {
	count, err := h.productUseCase.Count(c.Request.Context())
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