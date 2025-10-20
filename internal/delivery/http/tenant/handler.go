package tenant

import (
	"net/http"
	"strconv"

	tenantDomain "erp-api/internal/domain/tenant"
	tenantUseCase "erp-api/internal/usecase/tenant"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	useCase tenantUseCase.UseCaseInterface
}

func NewHandler(useCase tenantUseCase.UseCaseInterface) *Handler {
	return &Handler{useCase: useCase}
}

func (h *Handler) Create(c *gin.Context) {
	var dto tenantDomain.CreateTenantDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	tenant, err := h.useCase.Create(c.Request.Context(), &dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar tenant", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Tenant criado com sucesso", "data": tenant})
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID é obrigatório"})
		return
	}

	tenant, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tenant não encontrado", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenant encontrado", "data": tenant})
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID é obrigatório"})
		return
	}

	var dto tenantDomain.UpdateTenantDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	tenant, err := h.useCase.Update(c.Request.Context(), id, &dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar tenant", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenant atualizado com sucesso", "data": tenant})
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID é obrigatório"})
		return
	}

	err := h.useCase.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar tenant", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenant deletado com sucesso"})
}

func (h *Handler) List(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Limit inválido", "details": err.Error()})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Offset inválido", "details": err.Error()})
		return
	}

	tenants, err := h.useCase.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar tenants", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenants listados com sucesso", "data": tenants})
}

func (h *Handler) Count(c *gin.Context) {
	count, err := h.useCase.Count(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao contar tenants", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contagem realizada com sucesso", "data": gin.H{"count": count}})
}