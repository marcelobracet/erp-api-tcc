package settings

import (
	"net/http"

	settingsDomain "erp-api/internal/domain/settings"
	settingsUseCase "erp-api/internal/usecase/settings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	settingsUseCase settingsUseCase.UseCaseInterface
}

func NewHandler(settingsUseCase settingsUseCase.UseCaseInterface) *Handler {
	return &Handler{
		settingsUseCase: settingsUseCase,
	}
}

func (h *Handler) Get(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "tenant_id is required",
		})
		return
	}

	settings, err := h.settingsUseCase.Get(c.Request.Context(), tenantID)
	if err != nil {
		switch err {
		case settingsDomain.ErrSettingsNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Settings not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}
	
	c.JSON(http.StatusOK, settings)
}

func (h *Handler) Update(c *gin.Context) {
	var req settingsDomain.UpdateSettingsDTO
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	settings, err := h.settingsUseCase.Update(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case settingsDomain.ErrSettingsNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Settings not found",
			})
		case settingsDomain.ErrInvalidCNPJ:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid CNPJ",
			})
		case settingsDomain.ErrInvalidEmail:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid email",
			})
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		return
	}
	
	c.JSON(http.StatusOK, settings)
} 