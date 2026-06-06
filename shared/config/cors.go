package config

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS configuration constants.
const (
	envAllowedOrigins = "ALLOWED_ORIGINS"

	// devDefaultOrigins is used when ALLOWED_ORIGINS is empty in development.
	// These are the typical addresses for local front-end development servers
	// (Vite, CRA, Next.js dev server, etc.). They are safe for local use only
	// and must NEVER be trusted in production.
	devDefaultOrigins = "http://localhost:3000,http://localhost:5173,http://localhost:8080,http://127.0.0.1:3000,http://127.0.0.1:5173,http://127.0.0.1:8080"

	// prodDefaultOrigins is intentionally empty: a production deployment with
	// no explicit ALLOWED_ORIGINS must fail closed (no cross-origin access)
	// rather than silently allowing common local development hosts.
	prodDefaultOrigins = ""
)

// SetupCORS configures Cross-Origin Resource Sharing (CORS) for the given
// Gin router. It is designed for APIs that rely on session-based
// authentication, where the browser must be allowed to send cookies /
// credentials on cross-origin requests.
//
// The allowed origins are read from the ALLOWED_ORIGINS environment variable
// (comma-separated, e.g. "https://app.example.com,https://admin.example.com").
// When the variable is unset or empty, a sensible default is selected based on
// the provided environment:
//   - "development": a curated list of common local dev-server origins.
//   - anything else (including "production"): an empty list, which causes the
//     middleware to reject every cross-origin request. This is the safe
//     default: a misconfigured production server must not expose cookies to
//     arbitrary origins.
//
// IMPORTANT: When AllowCredentials is true, the CORS spec (and all major
// browsers) forbid using the wildcard "*" for Access-Control-Allow-Origin.
// The browser would refuse to attach cookies, and the request would fail.
// This middleware therefore reflects the exact request origin back to the
// client (only when that origin is in the allow-list) instead of using "*".
//
// Parameters:
//   - router: the *gin.Engine to attach the CORS middleware to. It is also
//     used to register a catch-all OPTIONS handler that returns 204 No
//     Content, satisfying CORS preflight requests.
//   - env: the current environment name. Pass "development" to enable the
//     development defaults; any other value is treated as production.
func SetupCORS(router *gin.Engine, env string) {
	isDevelopment := strings.EqualFold(env, "development")

	// Resolve the list of allowed origins.
	allowedOrigins := loadAllowedOrigins(isDevelopment)

	// Build a set for O(1) lookups.
	allowedSet := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		allowedSet[strings.TrimSpace(origin)] = struct{}{}
	}

	// Log the effective configuration so it is visible in the server boot logs.
	// We deliberately log the resolved list (not the raw env value) to make
	// debugging CORS issues easier: you can see exactly which origins will be
	// accepted without having to reproduce the env on your own machine.
	log.Printf("[CORS] Environment: %s", envLabel(env, isDevelopment))
	log.Printf("[CORS] AllowCredentials: true (required for session cookies)")
	log.Printf("[CORS] Allowed origins (%d): %s", len(allowedOrigins), formatOrigins(allowedOrigins))

	if !isDevelopment && len(allowedOrigins) == 0 {
		log.Printf("[CORS] WARNING: no allowed origins configured for production. " +
			"All cross-origin requests will be rejected. " +
			"Set ALLOWED_ORIGINS to a comma-separated list of trusted origins.")
	}

	router.Use(corsMiddleware(allowedSet, isDevelopment))

	// Register a catch-all OPTIONS handler that short-circuits preflight
	// requests with 204 No Content, as recommended by the CORS spec and by
	// Gin's documentation. The middleware above already sets the appropriate
	// Access-Control-* response headers, so all we need to do here is abort
	// the chain after writing the status code.
	router.OPTIONS("/*any", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNoContent)
	})
}

// loadAllowedOrigins parses the ALLOWED_ORIGINS environment variable and
// returns the resolved list of allowed origins, falling back to an
// environment-specific default when the variable is unset or empty.
func loadAllowedOrigins(isDevelopment bool) []string {
	raw := strings.TrimSpace(os.Getenv(envAllowedOrigins))

	if raw == "" {
		defaults := prodDefaultOrigins
		if isDevelopment {
			defaults = devDefaultOrigins
		}
		raw = defaults
	}

	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			origins = append(origins, p)
		}
	}
	return origins
}

// corsMiddleware returns a Gin handler that sets the CORS response headers
// on every request and short-circuits preflight (OPTIONS) requests with a
// 204 No Content status.
//
// Headers set:
//   - Access-Control-Allow-Origin:    the request origin, but only if it is
//     in the allow-list. A wildcard "*" is intentionally NOT used, because
//     AllowCredentials is true: the CORS spec forbids that combination, and
//     browsers would silently drop the session cookie. We must echo the
//     exact origin instead.
//   - Vary: Origin:                   tells shared caches and proxies that
//     the response varies per origin, preventing one user's CORS-enabled
//     response from being served to a different origin.
//   - Access-Control-Allow-Credentials: true
//   - Access-Control-Allow-Methods:    differs by environment; production
//     uses a minimal, explicit list, while development allows the full
//     common set for convenience.
//   - Access-Control-Allow-Headers:    Authorization, Content-Type, and
//     X-Session-ID (the header used by SessionAuth middleware). X-Session-ID
//     is required so the browser is allowed to send the session token as a
//     custom header in cross-origin requests.
//   - Access-Control-Max-Age:          how long the preflight result can be
//     cached by the browser. Longer in development to reduce preflight
//     chatter, shorter in production to limit the blast radius if allowed
//     methods/headers ever change.
func corsMiddleware(allowed map[string]struct{}, isDevelopment bool) gin.HandlerFunc {
	allowMethods := "GET,POST,PUT,PATCH,DELETE,OPTIONS"
	allowHeaders := "Origin,Content-Type,Accept,Authorization,X-Session-ID"
	maxAge := "600" // 10 minutes
	if isDevelopment {
		// Development gets a friendlier, more permissive set so live-reload
		// tools and quick experiments don't have to restart the server.
		allowMethods = "GET,POST,PUT,PATCH,DELETE,OPTIONS,HEAD"
		allowHeaders = "Origin,Content-Type,Accept,Authorization,X-Session-ID,X-Requested-With"
		maxAge = "86400" // 24 hours, to keep dev iteration fast
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Same-origin requests (no Origin header) are passed through with
		// no CORS headers; browsers don't enforce CORS for them anyway.
		if origin != "" {
			if _, ok := allowed[origin]; ok {
				// Echo the exact origin. We MUST NOT use "*" here because
				// AllowCredentials is true; see the package-level comment
				// in SetupCORS for the full explanation.
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Vary", "Origin")
				c.Header("Access-Control-Allow-Credentials", "true")
				c.Header("Access-Control-Allow-Methods", allowMethods)
				c.Header("Access-Control-Allow-Headers", allowHeaders)
				c.Header("Access-Control-Max-Age", maxAge)
			} else if !isDevelopment {
				// In production, log rejected origins so operators can spot
				// misconfigured clients. We still let the request continue;
				// the browser will block the response on the client side
				// because the Allow-Origin header is missing.
				log.Printf("[CORS] Rejected cross-origin request from origin: %s (method=%s path=%s)",
					origin, c.Request.Method, c.Request.URL.Path)
			}
		}

		// Short-circuit preflight requests. The catch-all OPTIONS route
		// registered in SetupCORS also handles this, but responding here as
		// well makes the middleware robust to direct handler invocations
		// and unit tests.
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// envLabel returns a human-readable label for the environment, used in logs.
func envLabel(env string, isDevelopment bool) string {
	if env == "" {
		if isDevelopment {
			return "development (inferred)"
		}
		return "production (inferred)"
	}
	if isDevelopment {
		return "development"
	}
	return "production"
}

// formatOrigins renders an origin slice in a compact, log-friendly form.
// It returns "<none>" for an empty slice so logs don't show a confusing
// trailing colon or empty parentheses.
func formatOrigins(origins []string) string {
	if len(origins) == 0 {
		return "<none>"
	}
	return strings.Join(origins, ", ")
}
