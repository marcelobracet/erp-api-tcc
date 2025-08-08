package client

import (
	"net/http"
	"strconv"

	clientDomain "erp-api/internal/domain/client"
	clientUseCase "erp-api/internal/usecase/client"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	clientUseCase clientUseCase.UseCaseInterface
}

func NewHandler(clientUseCase clientUseCase.UseCaseInterface) *Handler {
	return &Handler{
		clientUseCase: clientUseCase,
	}
}

// Create cria um novo cliente
func (h *Handler) Create(c *gin.Context) {
	var req clientDomain.CreateClientDTO
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	client, err := h.clientUseCase.Create(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case clientDomain.ErrClientAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{
				"error": "Client already exists",
			})
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		return
	}
	
	c.JSON(http.StatusCreated, client)
}

// GetByID obtém um cliente por ID
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	
	client, err := h.clientUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		switch err {
		case clientDomain.ErrClientNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Client not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}
	
	c.JSON(http.StatusOK, client)
}

// Update atualiza um cliente
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var req clientDomain.UpdateClientDTO
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	client, err := h.clientUseCase.Update(c.Request.Context(), id, &req)
	if err != nil {
		switch err {
		case clientDomain.ErrClientNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Client not found",
			})
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		return
	}
	
	c.JSON(http.StatusOK, client)
}

// Delete deleta um cliente
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	
	err := h.clientUseCase.Delete(c.Request.Context(), id)
	if err != nil {
		switch err {
		case clientDomain.ErrClientNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Client not found",
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

// List lista clientes
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
	
	clients, err := h.clientUseCase.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	total, err := h.clientUseCase.Count(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	response := clientDomain.ClientListDTO{
		Clients: make([]*clientDomain.ClientDTO, len(clients)),
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	}
	
	for i, client := range clients {
		response.Clients[i] = &clientDomain.ClientDTO{
			ID:           client.ID,
			Name:         client.Name,
			Email:        client.Email,
			Phone:        client.Phone,
			Document:     client.Document,
			DocumentType: client.DocumentType,
			Address:      client.Address,
			City:         client.City,
			State:        client.State,
			ZipCode:      client.ZipCode,
			IsActive:     client.IsActive,
			CreatedAt:    client.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:    client.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	
	c.JSON(http.StatusOK, response)
}

// Count conta clientes
func (h *Handler) Count(c *gin.Context) {
	count, err := h.clientUseCase.Count(c.Request.Context())
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