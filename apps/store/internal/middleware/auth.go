package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// APIKeyAuth creates a middleware that validates API key authentication.
// If apiKey is empty, the middleware is a no-op (authentication disabled).
// When enabled, requests must include an "X-API-Key" header matching the configured key.
// Public endpoints like /health and /docs are always excluded from authentication.
func APIKeyAuth(api huma.API, apiKey string) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		// Skip auth for public endpoints
		path := ctx.URL().Path
		if path == "/health" || path == "/docs" || path == "/openapi.json" || path == "/openapi.yaml" {
			next(ctx)
			return
		}

		// Get API key from header
		providedKey := ctx.Header("X-API-Key")

		// Validate API key using constant-time comparison to prevent timing attacks
		if !secureCompare(providedKey, apiKey) {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized,
				"Invalid or missing API key",
			)
			return
		}

		// Continue to next middleware/handler
		next(ctx)
	}
}

// secureCompare performs a constant-time comparison of two strings to prevent timing attacks.
// Returns true if the strings are equal, false otherwise.
func secureCompare(a, b string) bool {
	// If lengths don't match, still compare to prevent timing leaks
	aBytes := []byte(a)
	bBytes := []byte(b)

	// subtle.ConstantTimeCompare requires equal length, so we need to handle length mismatch
	if len(aBytes) != len(bBytes) {
		return false
	}

	return subtle.ConstantTimeCompare(aBytes, bBytes) == 1
}
