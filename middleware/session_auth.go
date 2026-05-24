package middleware

import (
	"AnbariAPI/dto"
	"AnbariAPI/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SessionAuth(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionIDStr := c.GetHeader("X-Session-ID")
		if sessionIDStr == "" {
			cookie, err := c.Cookie("session_id")
			if err != nil || cookie == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "session_id required"})
				c.Abort()
				return
			}
			sessionIDStr = cookie
		}

		sessionID, err := uuid.Parse(sessionIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session_id"})
			c.Abort()
			return
		}

		user, err := authService.ValidateSession(sessionID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired session"})
			c.Abort()
			return
		}

		currentUser := dto.UserDTO{
			ID:        user.ID,
			Email:     user.Email,
			Phone:     user.Phone,
			CreatedAt: user.CreatedAt,
		}

		c.Set("currentUser", &currentUser)
		c.Set("sessionID", sessionID)
		c.Next()
	}
}