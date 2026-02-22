package user

import (
	"net/http"
	"os"
	"strconv"

	userDomain "erp-api/internal/domain/user"
	userUseCase "erp-api/internal/usecase/user"
	"erp-api/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// Handler implementa os HTTP handlers para usuários
type Handler struct {
	userUseCase userUseCase.UseCaseInterface
}

// NewHandler cria uma nova instância do Handler
func NewHandler(userUseCase userUseCase.UseCaseInterface) *Handler {
	return &Handler{
		userUseCase: userUseCase,
	}
}

// Register registra um novo usuário
// @Summary Registrar novo usuário
// @Description Registra um novo usuário no sistema
// @Tags users
// @Accept json
// @Produce json
// @Param user body userDomain.CreateUserRequest true "Dados do usuário"
// @Success 201 {object} userDomain.User
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /users/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req userDomain.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	user, err := h.userUseCase.Register(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case userDomain.ErrUserAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{
				"error": "User already exists",
			})
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetByID busca um usuário pelo ID
// @Summary Buscar usuário por ID
// @Description Busca um usuário específico pelo ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Success 200 {object} userDomain.User
// @Failure 404 {object} map[string]interface{}
// @Router /users/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		switch err {
		case userDomain.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetProfile retorna o perfil do usuário autenticado
// @Summary Obter perfil do usuário
// @Description Retorna o perfil do usuário autenticado
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} userDomain.User
// @Failure 401 {object} map[string]interface{}
// @Router /users/profile [get]
func (h *Handler) GetProfile(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	user, err := h.userUseCase.GetByID(c.Request.Context(), userID)
	if err != nil {
		if err == userDomain.ErrUserNotFound && os.Getenv("AUTH_PROVIDER") == "keycloak" {
			tenantID, _ := middleware.GetTenantIDFromContext(c)
			emailValue, hasEmail := middleware.GetUserEmailFromContext(c)
			var email *string
			if hasEmail && emailValue != "" {
				email = &emailValue
			}

			displayName := userID
			if email != nil {
				displayName = *email
			}

			_, _ = h.userUseCase.Register(c.Request.Context(), &userDomain.CreateUserRequest{
				TenantID:    tenantID,
				KeycloakID:  userID,
				DisplayName: displayName,
				Email:       email,
			})

			user, err = h.userUseCase.GetByID(c.Request.Context(), userID)
			if err == nil {
				c.JSON(http.StatusOK, user)
				return
			}
		}

		switch err {
		case userDomain.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// Update atualiza um usuário
// @Summary Atualizar usuário
// @Description Atualiza os dados de um usuário
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Param user body userDomain.UpdateUserRequest true "Dados para atualização"
// @Success 200 {object} userDomain.User
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	var req userDomain.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	user, err := h.userUseCase.Update(c.Request.Context(), id, &req)
	if err != nil {
		switch err {
		case userDomain.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// Delete remove um usuário
// @Summary Deletar usuário
// @Description Remove um usuário do sistema
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]interface{}
// @Router /users/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.userUseCase.Delete(c.Request.Context(), id)
	if err != nil {
		switch err {
		case userDomain.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// List retorna uma lista de usuários
// @Summary Listar usuários
// @Description Retorna uma lista paginada de usuários
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int false "Limite de resultados (padrão: 10)"
// @Param offset query int false "Offset para paginação (padrão: 0)"
// @Success 200 {array} userDomain.User
// @Router /users [get]
func (h *Handler) List(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	users, err := h.userUseCase.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, users)
}

// Count retorna o total de usuários
// @Summary Contar usuários
// @Description Retorna o total de usuários no sistema
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]int
// @Router /users/count [get]
func (h *Handler) Count(c *gin.Context) {
	count, err := h.userUseCase.Count(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": count,
	})
}
