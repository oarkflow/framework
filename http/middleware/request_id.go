package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sujit-baniya/framework/contracts/http"
	"github.com/sujit-baniya/framework/utils/xid"
)

// ConfigRequestID defines the config for middleware.
type ConfigRequestID struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c http.Context) bool

	// Header is the header key where to get/set the unique request ID
	//
	// Optional. Default: "X-Request-ID"
	Header string

	// Generator defines a function to generate the unique identifier.
	//
	// Optional. Default: utils.UUID
	Generator func() string

	// ContextKey defines the key used when storing the request ID in
	// the locals for a specific request.
	//
	// Optional. Default: requestid
	ContextKey string
}

// ConfigRequestIDDefault is the default config
var ConfigRequestIDDefault = ConfigRequestID{
	Next:       nil,
	Header:     fiber.HeaderXRequestID,
	Generator:  xid.New().String,
	ContextKey: "requestid",
}

// Helper function to set default values
func configRequestIDDefault(config ...ConfigRequestID) ConfigRequestID {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigRequestIDDefault
	}

	// Override default config
	cfg := config[0]

	// Set default values
	if cfg.Header == "" {
		cfg.Header = ConfigRequestIDDefault.Header
	}
	if cfg.Generator == nil {
		cfg.Generator = ConfigRequestIDDefault.Generator
	}
	if cfg.ContextKey == "" {
		cfg.ContextKey = ConfigRequestIDDefault.ContextKey
	}
	return cfg
}

// RequestID creates a new middleware handler
func RequestID(config ...ConfigRequestID) http.HandlerFunc {
	// Set default config
	cfg := configRequestIDDefault(config...)

	// Return new handler
	return func(c http.Context) error {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}
		// Get id from request, else we generate one
		rid := c.Header(cfg.Header, cfg.Generator())

		// Set new id to response header
		c.SetHeader(cfg.Header, rid)

		// Add the request ID to locals
		c.WithValue(cfg.ContextKey, rid)

		// Continue stack
		return c.Next()
	}
}
