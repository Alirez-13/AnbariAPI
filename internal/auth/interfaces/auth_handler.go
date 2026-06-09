// # SINGLE REASON: Handle auth HTTP requests.
package interfaces

import (
	"net/http"

	"AnbariAPI/internal/auth/application"
	"AnbariAPI/internal/auth/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SessionManager interface {
	Create(c *gin.Context, sessionID string) error
	Destroy(c *gin.Context) error
	SessionID(c *gin.Context) (string, error)
}

type AuthHandler struct {
	authService    *application.AuthService
	sessionManager SessionManager
}

func NewAuthHandler(authService *application.AuthService, sessionManager SessionManager) *AuthHandler {
	return &AuthHandler{authService: authService, sessionManager: sessionManager}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := h.authService.Register(application.RegisterInput(req))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.sessionManager.Create(c, response.SessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session cookie"})
		return
	}
	c.JSON(http.StatusCreated, authResponse(response))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := h.authService.Login(application.LoginInput(req))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if err := h.sessionManager.Create(c, response.SessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session cookie"})
		return
	}
	c.JSON(http.StatusOK, authResponse(response))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	sessionIDStr, err := h.sessionManager.SessionID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session cookie"})
		return
	}
	if sessionIDStr == "" {
		sessionIDStr = c.GetHeader("X-Session-ID")
	}
	if sessionIDStr == "" {
		sessionIDStr = c.Query("session_id")
	}
	if sessionIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id required"})
		return
	}
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session_id"})
		return
	}
	if err := h.authService.Logout(sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout"})
		return
	}
	if err := h.sessionManager.Destroy(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to destroy session cookie"})
		return
	}
	c.JSON(http.StatusOK, dto.LogoutResponse{Message: "logged out successfully"})
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.JSON(http.StatusOK, user.(*dto.UserDTO))
}

func authResponse(result *application.AuthResult) *dto.AuthResponse {
	return &dto.AuthResponse{
		User:      dto.UserDTO{ID: result.User.ID, Email: result.User.Email, Phone: result.User.Phone, CreatedAt: result.User.CreatedAt},
		SessionID: result.SessionID,
		ExpiresAt: result.ExpiresAt,
	}
}
