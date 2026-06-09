// # SINGLE REASON: Manage auth session cookies and Gin session middleware.
package interfaces

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"AnbariAPI/internal/auth/application"
	"AnbariAPI/internal/auth/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

const (
	envSessionDomain  = "SESSION_DOMAIN"
	sessionCookieName = "session_id"
	sessionValueKey   = "session_id"
)

type CookieSessionManager struct {
	store  *sessions.CookieStore
	logger *slog.Logger
}

func NewSessionManager(secret, env string, logger *slog.Logger) *CookieSessionManager {
	if logger == nil {
		logger = slog.Default()
	}
	store := sessions.NewCookieStore([]byte(secret))
	store.Options = &sessions.Options{
		Path:     "/",
		Domain:   strings.TrimSpace(os.Getenv(envSessionDomain)),
		MaxAge:   int(application.SessionDuration.Seconds()),
		HttpOnly: true,
		Secure:   !strings.EqualFold(env, "development"),
		SameSite: http.SameSiteLaxMode,
	}
	return &CookieSessionManager{store: store, logger: logger}
}

func (m *CookieSessionManager) Create(c *gin.Context, sessionID string) error {
	session, err := m.store.New(c.Request, sessionCookieName)
	if err != nil {
		return err
	}
	session.Values[sessionValueKey] = sessionID
	session.Options.MaxAge = int(application.SessionDuration.Seconds())
	if err := session.Save(c.Request, c.Writer); err != nil {
		return err
	}
	m.logger.Info("session created", slog.String("session_id", sessionID))
	return nil
}

func (m *CookieSessionManager) Refresh(c *gin.Context) error {
	session, err := m.store.Get(c.Request, sessionCookieName)
	if err != nil {
		return err
	}
	if _, ok := session.Values[sessionValueKey]; !ok {
		return nil
	}
	session.Options.MaxAge = int(application.SessionDuration.Seconds())
	return session.Save(c.Request, c.Writer)
}

func (m *CookieSessionManager) SessionID(c *gin.Context) (string, error) {
	session, err := m.store.Get(c.Request, sessionCookieName)
	if err != nil {
		return "", err
	}
	sessionID, _ := session.Values[sessionValueKey].(string)
	return sessionID, nil
}

func (m *CookieSessionManager) Destroy(c *gin.Context) error {
	session, err := m.store.Get(c.Request, sessionCookieName)
	if err != nil {
		return err
	}
	sessionID, _ := session.Values[sessionValueKey].(string)
	delete(session.Values, sessionValueKey)
	session.Options.MaxAge = -1
	if err := session.Save(c.Request, c.Writer); err != nil {
		return err
	}
	m.logger.Info("session destroyed", slog.String("session_id", sessionID))
	return nil
}

func SessionAuth(authService *application.AuthService, sessionManager *CookieSessionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionIDStr, err := sessionManager.SessionID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session cookie"})
			c.Abort()
			return
		}
		if sessionIDStr == "" {
			sessionIDStr = c.GetHeader("X-Session-ID")
			if sessionIDStr == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "session_id required"})
				c.Abort()
				return
			}
		}

		sessionID, err := uuid.Parse(sessionIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session_id"})
			c.Abort()
			return
		}

		user, err := authService.ValidateSession(sessionID)
		if err != nil {
			if errors.Is(err, application.ErrSessionExpired) {
				if destroyErr := sessionManager.Destroy(c); destroyErr != nil {
					slog.Warn("failed to destroy expired session cookie", slog.Any("error", destroyErr))
				}
				c.JSON(http.StatusUnauthorized, gin.H{"error": "session expired"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
			}
			c.Abort()
			return
		}
		if err := sessionManager.Refresh(c); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to refresh session"})
			c.Abort()
			return
		}

		currentUser := dto.UserDTO{ID: user.ID, Email: user.Email, Phone: user.Phone, CreatedAt: user.CreatedAt}
		c.Set("currentUser", &currentUser)
		c.Set("sessionID", sessionID)
		c.Next()
	}
}
