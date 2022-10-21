package middleware

import (
	"fmt"
	"github.com/sujit-baniya/framework/contracts/http"
	"os"
	"runtime"
)

// ConfigRecover defines the config for middleware.
type ConfigRecover struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c http.Context) bool

	// EnableStackTrace enables handling stack trace
	//
	// Optional. Default: false
	EnableStackTrace bool

	// StackTraceHandler defines a function to handle stack trace
	//
	// Optional. Default: defaultStackTraceHandler
	StackTraceHandler func(c http.Context, e interface{})
}

var defaultStackTraceBufLen = 1024

// ConfigRecoverDefault is the default config
var ConfigRecoverDefault = ConfigRecover{
	Next:              nil,
	EnableStackTrace:  false,
	StackTraceHandler: defaultStackTraceHandler,
}

// Helper function to set default values
func configRecoverDefault(config ...ConfigRecover) ConfigRecover {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigRecoverDefault
	}

	// Override default config
	cfg := config[0]

	if cfg.EnableStackTrace && cfg.StackTraceHandler == nil {
		cfg.StackTraceHandler = defaultStackTraceHandler
	}

	return cfg
}

func defaultStackTraceHandler(c http.Context, e interface{}) {
	buf := make([]byte, defaultStackTraceBufLen)
	buf = buf[:runtime.Stack(buf, false)]
	_, _ = os.Stderr.WriteString(fmt.Sprintf("panic: %v\n%s\n", e, buf))
}

// Recover creates a new middleware handler
func Recover(config ...ConfigRecover) http.HandlerFunc {
	// Set default config
	cfg := configRecoverDefault(config...)

	// Return new handler
	return func(c http.Context) (err error) {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Catch panics
		defer func() {
			if r := recover(); r != nil {
				if cfg.EnableStackTrace {
					cfg.StackTraceHandler(c, r)
				}

				var ok bool
				if err, ok = r.(error); !ok {
					// Set error that will call the global error handler
					err = fmt.Errorf("%v", r)
				}
			}
		}()

		// Return err if existed, else move to next handler
		return c.Next()
	}
}
