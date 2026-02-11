package client

import (
	"net/http"
	"strconv"

	clientDomain "erp-api/internal/domain/client"
	clientUseCase "erp-api/internal/usecase/client"
	"erp-api/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
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
	log.Info().Msg("Create client started")

	tenantID, exists := middleware.GetTenantIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var req clientDomain.CreateClientDTO

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

	log.Info().Msg("Create client ended")
	c.JSON(http.StatusCreated, client)
}

// GetByID obtém um cliente por ID
func (h *Handler) GetByID(c *gin.Context) {
	log.Info().Msg("Get client by ID started")
	tenantID, exists := middleware.GetTenantIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	id := c.Param("id")

	client, err := h.clientUseCase.GetByID(c.Request.Context(), tenantID, id)
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

	log.Info().Msg("Get client by ID ended")
	c.JSON(http.StatusOK, client)
}

// Update atualiza um cliente
func (h *Handler) Update(c *gin.Context) {
	log.Info().Msg("Update client started")

	tenantID, exists := middleware.GetTenantIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	id := c.Param("id")
	var req clientDomain.UpdateClientDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	client, err := h.clientUseCase.Update(c.Request.Context(), tenantID, id, &req)
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

	log.Info().Msg("Update client ended")
	c.JSON(http.StatusOK, client)
}

// Delete deleta um cliente
func (h *Handler) Delete(c *gin.Context) {
	log.Info().Msg("Delete client started")

	tenantID, exists := middleware.GetTenantIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	id := c.Param("id")

	err := h.clientUseCase.Delete(c.Request.Context(), tenantID, id)
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

	log.Info().Msg("Delete client ended")
	c.Status(http.StatusNoContent)
}

// List lista clientes
func (h *Handler) List(c *gin.Context) {
	log.Info().Msg("List clients started")

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

	clients, err := h.clientUseCase.List(c.Request.Context(), tenantID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	total, err := h.clientUseCase.Count(c.Request.Context(), tenantID)
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

	log.Info().Msg("List clients ended")
	c.JSON(http.StatusOK, response)
}

// Count conta clientes
func (h *Handler) Count(c *gin.Context) {
	log.Info().Msg("Count clients started")

	tenantID, exists := middleware.GetTenantIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	count, err := h.clientUseCase.Count(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Info().Msg("Count clients ended")
	c.JSON(http.StatusOK, gin.H{
		"count": count,
	})
}
