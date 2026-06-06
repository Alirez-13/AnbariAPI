package main

import (
	"AnbariAPI/Internal/auth/handler"
	"AnbariAPI/Internal/auth/middleware"
	"AnbariAPI/Internal/auth/repository"
	"AnbariAPI/Internal/auth/service"
	"AnbariAPI/api/routes"
	corscfg "AnbariAPI/shared/config"
	"AnbariAPI/shared/database"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	envPort         = "PORT"
	envAppEnv       = "APP_ENV"
	envSessionKey   = "SESSION_SECRET"
	defaultPort     = "8080"
	defaultEnv      = "development"
	readTimeout     = 15 * time.Second
	writeTimeout    = 15 * time.Second
	idleTimeout     = 60 * time.Second
	shutdownTimeout = 30 * time.Second
)

func main() {
	// Structured JSON logger. Using slog as the single source of truth for
	// logs makes it easy to ship to a log aggregator later (Loki, ELK, etc.)
	// without reformatting. We install it as the slog default so every
	// dependency that uses slog.Default() (e.g. the auth service) also gets
	// JSON output.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	if err := run(logger); err != nil {
		logger.Error("server exited with error", slog.Any("error", err))
		os.Exit(1)
	}
}

// run is split out of main so the entire bootstrap (config loading, DB init,
// server lifecycle) is testable and so all errors funnel through one place
// rather than calling log.Fatal from inside helper functions.
func run(logger *slog.Logger) error {
	env := getEnv(envAppEnv, defaultEnv)
	port := getEnv(envPort, defaultPort)

	// Session secret: a random secret is required in production (otherwise
	// signed cookies can be forged). In development, we generate one per
	// process so the app boots without configuration, but we log a warning
	// so it's obvious during local development.
	sessionSecret, err := loadSessionSecret(env, logger)
	if err != nil {
		return err
	}
	// Database. DefaultConfig() is overridden by the same env-driven path
	// the rest of the project uses; we keep the inline path here for
	// consistency with the existing bootstrap until a unified config
	// package exists.
	dbConfig := database.DefaultConfig()
	dbConfig.DBPath = getEnv("DB_PATH", "./data/inventory.db")
	dbConfig.Debug = env == "development"

	if err := database.InitDatabase(dbConfig); err != nil {
		return err
	}
	defer func() {
		if cerr := database.CloseDatabase(); cerr != nil {
			logger.Error("failed to close database", slog.Any("error", cerr))
		}
	}()

	db := database.GetDB()

	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	authService := service.NewAuthService(userRepo, sessionRepo, logger)
	sessionManager := middleware.NewSessionManager(sessionSecret, env, logger)
	authHandler := handler.NewAuthHandler(authService, sessionManager)

	// Build the router with the required middleware order:
	//   1. Recovery  — catches panics, returns 500 instead of crashing
	//   2. Logger    — request/response logging (Gin's built-in)
	//   3. CORS      — must run after Logger so preflight rejections are logged
	//   4. Routes    — application handlers; per-group session auth is added
	//                  below on the protected group, not globally, because
	//                  /auth/* and /health are intentionally anonymous.
	// We deliberately don't use gin.Default() so the order is explicit, and
	// we control exactly which middlewares run in which sequence.
	gin.SetMode(gin.ReleaseMode)
	if env == "development" {
		gin.SetMode(gin.DebugMode)
	}
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			// Hand off to slog so all logs are JSON-shaped and include
			// the standard request fields. A failed request (status >= 500
			// or any error from a handler) is logged at Error level.
			level := slog.LevelInfo
			if param.StatusCode >= http.StatusInternalServerError {
				level = slog.LevelError
			}
			logger.LogAttrs(context.Background(), level, "http request",
				slog.String("method", param.Method),
				slog.String("path", param.Path),
				slog.Int("status", param.StatusCode),
				slog.Duration("latency", param.Latency),
				slog.String("client_ip", param.ClientIP),
				slog.String("user_agent", param.Request.UserAgent()),
				slog.String("errors", param.ErrorMessage),
			)
			return ""
		},
		Output: os.Stdout,
	}))

	// CORS is installed BEFORE route registration so it can register its
	// catch-all OPTIONS handler with priority over user-defined routes.
	corscfg.SetupCORS(router, env)

	// Health check. Registered before the protected group so liveness
	// probes never require a session.
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Application routes (categories, products, inventory, auth).
	routes.SetupRoutes(router, db, authHandler)

	// Protected routes: per-group session auth, in addition to the global
	// Recovery + Logger + CORS chain. The order is therefore
	// Recovery → Logger → CORS → SessionAuth → handler.
	protected := router.Group("/api")
	protected.Use(middleware.SessionAuth(authService, sessionManager))
	{
		protected.GET("/me", authHandler.GetCurrentUser)
	}

	// http.Server (instead of r.Run) is required for two reasons:
	//   1. Timeouts: r.Run uses zero timeouts, which leaves the server open
	//      to slow-loris attacks and runaway handlers.
	//   2. Graceful shutdown: we need a handle on the *http.Server to call
	//      Shutdown() when a signal arrives.
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	// signal.NotifyContext gives us a context that is canceled on
	// SIGINT/SIGTERM, which is the idiomatic Go 1.16+ way to wait for
	// shutdown signals. Using context (not a channel) means the cancel
	// function is wired up for us and we can pass it down to shut down.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		logger.Info("server starting",
			slog.String("addr", srv.Addr),
			slog.String("env", env),
		)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		logger.Info("shutdown signal received, draining connections",
			slog.Duration("timeout", shutdownTimeout),
		)
	}

	// Shutdown context bounded by shutdownTimeout. We give in-flight
	// requests up to that long to finish before we force-close the
	// server and exit.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}
	logger.Info("server stopped cleanly")
	return nil
}

// getEnv returns the value of the named env var, or fallback if unset/empty.
func getEnv(name, fallback string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return fallback
}

// loadSessionSecret returns SESSION_SECRET from the environment, falling
// back to a per-process random key in development. A persistent secret is
// mandatory in production; if it's missing we fail at boot rather than
// silently using an ephemeral key (which would invalidate every session
// on restart and mask the real configuration problem).
func loadSessionSecret(env string, logger *slog.Logger) (string, error) {
	if secret := os.Getenv(envSessionKey); secret != "" {
		return secret, nil
	}
	if env != "development" {
		return "", errors.New("SESSION_SECRET is required in non-development environments")
	}
	secret, err := randomHex(32)
	if err != nil {
		return "", err
	}
	logger.Warn("SESSION_SECRET not set; generated an ephemeral key for this process",
		slog.String("env", env),
		slog.String("note", "all sessions will be invalidated on restart"),
	)
	return secret, nil
}

// randomHex returns n random bytes encoded as hex. Used to derive the
// per-process development session secret.
func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
