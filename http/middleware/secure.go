package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sujit-baniya/framework/contracts/http"
)

// ConfigSecure ...
type ConfigSecure struct {
	// Filter defines a function to skip middleware.
	// Optional. Default: nil
	Filter func(http.Context) bool
	// XSSProtection
	// Optional. Default value "1; mode=block".
	XSSProtection string
	// ContentTypeNosniff
	// Optional. Default value "nosniff".
	ContentTypeNosniff string
	// XFrameOptions
	// Optional. Default value "SAMEORIGIN".
	// Possible values: "SAMEORIGIN", "DENY", "ALLOW-FROM uri"
	XFrameOptions string
	// HSTSMaxAge
	// Optional. Default value 0.
	HSTSMaxAge int
	// HSTSExcludeSubdomains
	// Optional. Default value false.
	HSTSExcludeSubdomains bool
	// ContentSecurityPolicy
	// Optional. Default value "".
	ContentSecurityPolicy string
	// CSPReportOnly
	// Optional. Default value false.
	CSPReportOnly bool
	// HSTSPreloadEnabled
	// Optional.  Default value false.
	HSTSPreloadEnabled bool
	// ReferrerPolicy
	// Optional. Default value "".
	ReferrerPolicy string

	// Permissions-Policy
	// Optional. Default value "".
	PermissionPolicy string
}

// Secure ...
func Secure(config ...ConfigSecure) http.HandlerFunc {
	// Init config
	var cfg ConfigSecure
	if len(config) > 0 {
		cfg = config[0]
	}
	// Set config default values
	if cfg.XSSProtection == "" {
		cfg.XSSProtection = "1; mode=block"
	}
	if cfg.ContentTypeNosniff == "" {
		cfg.ContentTypeNosniff = "nosniff"
	}
	if cfg.XFrameOptions == "" {
		cfg.XFrameOptions = "SAMEORIGIN"
	}
	// Return middleware handler
	return func(c http.Context) error {
		// Filter request to skip middleware
		if cfg.Filter != nil && cfg.Filter(c) {
			return c.Next()
		}
		if cfg.XSSProtection != "" {
			c.SetHeader(fiber.HeaderXXSSProtection, cfg.XSSProtection)
		}
		if cfg.ContentTypeNosniff != "" {
			c.SetHeader(fiber.HeaderXContentTypeOptions, cfg.ContentTypeNosniff)
		}
		if cfg.XFrameOptions != "" {
			c.SetHeader(fiber.HeaderXFrameOptions, cfg.XFrameOptions)
		}
		//@TODO - Add Secure() only after Gin TLS is identified
		// if (c.Secure() || (c.Header(fiber.HeaderXForwardedProto, "") == "https")) && cfg.HSTSMaxAge != 0 {
		if (c.Header(fiber.HeaderXForwardedProto, "") == "https") && cfg.HSTSMaxAge != 0 {
			subdomains := ""
			if !cfg.HSTSExcludeSubdomains {
				subdomains = "; includeSubdomains"
			}
			if cfg.HSTSPreloadEnabled {
				subdomains = fmt.Sprintf("%s; preload", subdomains)
			}
			c.SetHeader(fiber.HeaderStrictTransportSecurity, fmt.Sprintf("max-age=%d%s", cfg.HSTSMaxAge, subdomains))
		}
		if cfg.ContentSecurityPolicy != "" {
			if cfg.CSPReportOnly {
				c.SetHeader(fiber.HeaderContentSecurityPolicyReportOnly, cfg.ContentSecurityPolicy)
			} else {
				c.SetHeader(fiber.HeaderContentSecurityPolicy, cfg.ContentSecurityPolicy)
			}
		}
		if cfg.ReferrerPolicy != "" {
			c.SetHeader(fiber.HeaderReferrerPolicy, cfg.ReferrerPolicy)
		}
		if cfg.PermissionPolicy != "" {
			c.SetHeader(fiber.HeaderPermissionsPolicy, cfg.PermissionPolicy)

		}
		return c.Next()
	}
}
