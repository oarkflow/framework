package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sujit-baniya/framework/contracts/http"
	"github.com/sujit-baniya/framework/utils/xid"
	"strings"
	"time"

	"github.com/sujit-baniya/log"
)

type ConfigLog struct {
	Logger    *log.Logger
	LogWriter log.Writer
	RequestID func() string
}

// Log Middleware request_id + logger + recover for request traceability
func Log(config ConfigLog) http.HandlerFunc {
	return func(c http.Context) error {
		start := time.Now()
		if strings.Contains(c.Path(), "favicon") {
			return c.Next()
		}
		rid := c.Header(fiber.HeaderXRequestID, "")
		if config.RequestID == nil {
			config.RequestID = func() string {
				return xid.New().String()
			}
		}
		if rid == "" {
			rid = config.RequestID()
			c.SetHeader(fiber.HeaderXRequestID, rid)
		}
		nextHandler := c.Next()
		if c.Path() == "/" && c.Path() != c.Path() {
			return nextHandler
		}

		if config.Logger == nil {
			config.Logger = &log.Logger{
				TimeField:  "timestamp",
				TimeFormat: "2006-01-02 15:04:05",
			}
		}
		if config.LogWriter != nil {
			config.Logger.Writer = config.LogWriter
		}
		ip := c.Ip()
		curIP := c.Value("ip")
		if curIP != nil {
			ip = curIP.(string)
		}
		logging := log.NewContext(nil).
			Str("request_id", rid).
			Str("remote_ip", ip).
			Str("method", c.Method()).
			Str("host", c.Origin().Host).
			Str("path", c.Path()).
			Str("protocol", c.Origin().Proto).
			Int("status", c.StatusCode()).
			Str("latency", fmt.Sprintf("%s", time.Since(start))).
			Str("ua", c.Header(fiber.HeaderUserAgent, ""))

		log.Info().Str("request_id", rid).
			Str("remote_ip", ip).
			Str("method", c.Method()).
			Str("host", c.Origin().Host).
			Str("path", c.Path()).
			Str("protocol", c.Origin().Proto).
			Int("status", c.StatusCode()).
			Str("latency", fmt.Sprintf("%s", time.Since(start))).
			Str("ua", c.Header(fiber.HeaderUserAgent, ""))

		if nextHandler != nil {
			log.Info().Str("error", nextHandler.Error())
			logging.Str("error", nextHandler.Error())
		}

		ctx := logging.Value()
		switch {
		case c.StatusCode() >= 500:
			config.Logger.Error().Context(ctx).Msg("server error")
			log.Error().Context(ctx).Msg("server error")
		case c.StatusCode() >= 400:
			config.Logger.Error().Context(ctx).Msg("client error")
			log.Error().Context(ctx).Msg("client error")
		case c.StatusCode() >= 300:
			config.Logger.Warn().Context(ctx).Msg("redirect")
			log.Info().Context(ctx).Msg("redirect")
		case c.StatusCode() >= 200:
			config.Logger.Info().Context(ctx).Msg("success")
			log.Info().Context(ctx).Msg("success")
		case c.StatusCode() >= 100:
			config.Logger.Info().Context(ctx).Msg("informative")
			log.Info().Context(ctx).Msg("informative")
		default:
			config.Logger.Warn().Context(ctx).Msg("unknown status")
			log.Info().Context(ctx).Msg("unknown status")
		}
		return nextHandler
	}
}
