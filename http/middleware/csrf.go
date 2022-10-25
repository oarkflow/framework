package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/sujit-baniya/framework/contracts/http"
	"github.com/sujit-baniya/framework/storage"
	"github.com/sujit-baniya/framework/utils"
	"github.com/sujit-baniya/framework/utils/msgp"
	"github.com/sujit-baniya/framework/utils/xid"
	"net/textproto"
	"strings"
	"sync"
	"time"
)

// ConfigCsrf defines the config for middleware.
type ConfigCsrf struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c http.Context) bool

	// KeyLookup is a string in the form of "<source>:<key>" that is used
	// to create an Extractor that extracts the token from the request.
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "param:<name>"
	// - "form:<name>"
	// - "cookie:<name>"
	//
	// Ignored if an Extractor is explicitly set.
	//
	// Optional. Default: "header:X-CSRF-Token"
	KeyLookup string

	// Name of the session cookie. This cookie will store session key.
	// Optional. Default value "csrf_".
	CookieName string

	// Domain of the CSRF cookie.
	// Optional. Default value "".
	CookieDomain string

	// Path of the CSRF cookie.
	// Optional. Default value "".
	CookiePath string

	// Indicates if CSRF cookie is secure.
	// Optional. Default value false.
	CookieSecure bool

	// Indicates if CSRF cookie is HTTP only.
	// Optional. Default value false.
	CookieHTTPOnly bool

	// Value of SameSite cookie.
	// Optional. Default value "Lax".
	CookieSameSite string

	// Decides whether cookie should last for only the browser sesison.
	// Ignores Expiration if set to true
	CookieSessionOnly bool

	// Expiration is the duration before csrf token will expire
	//
	// Optional. Default: 1 * time.Hour
	Expiration time.Duration

	// Store is used to store the state of the middleware
	//
	// Optional. Default: memory.New()
	Storage fiber.Storage

	// Context key to store generated CSRF token into context.
	// If left empty, token will not be stored in context.
	//
	// Optional. Default: ""
	ContextKey string

	// KeyGenerator creates a new CSRF token
	//
	// Optional. Default: utils.UUID
	KeyGenerator func() string

	// ErrorHandler is executed when an error is returned from fiber.Handler.
	//
	// Optional. Default: DefaultErrorHandler
	ErrorHandler http.ErrorHandler

	// Extractor returns the csrf token
	//
	// If set this will be used in place of an Extractor based on KeyLookup.
	//
	// Optional. Default will create an Extractor based on KeyLookup.
	Extractor func(c http.Context) (string, error)
}

const HeaderName = "X-Csrf-Token"

// ConfigCsrfDefault is the default config
var ConfigCsrfDefault = ConfigCsrf{
	KeyLookup:      "header:" + HeaderName,
	CookieName:     "csrf_token",
	CookieSameSite: "Lax",
	Expiration:     1 * time.Hour,
	KeyGenerator:   xid.New().String,
	ErrorHandler:   defaultErrorHandler,
	Extractor:      CsrfFromHeader(HeaderName),
}

// default ErrorHandler that process return error from fiber.Handler
var defaultErrorHandler = func(c http.Context, err error) error {
	c.AbortWithStatus(fiber.StatusForbidden)
	return fiber.ErrForbidden
}

// Helper function to set default values
func configDefault(config ...ConfigCsrf) ConfigCsrf {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigCsrfDefault
	}

	// Override default config
	cfg := config[0]
	if cfg.KeyLookup == "" {
		cfg.KeyLookup = ConfigCsrfDefault.KeyLookup
	}
	if int(cfg.Expiration.Seconds()) <= 0 {
		cfg.Expiration = ConfigCsrfDefault.Expiration
	}
	if cfg.CookieName == "" {
		cfg.CookieName = ConfigCsrfDefault.CookieName
	}
	if cfg.CookieSameSite == "" {
		cfg.CookieSameSite = ConfigCsrfDefault.CookieSameSite
	}
	if cfg.KeyGenerator == nil {
		cfg.KeyGenerator = ConfigCsrfDefault.KeyGenerator
	}
	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = ConfigCsrfDefault.ErrorHandler
	}

	// Generate the correct extractor to get the token from the correct location
	selectors := strings.Split(cfg.KeyLookup, ":")

	if len(selectors) != 2 {
		panic("[CSRF] KeyLookup must in the form of <source>:<key>")
	}

	if cfg.Extractor == nil {
		// By default, we extract from a header
		cfg.Extractor = CsrfFromHeader(textproto.CanonicalMIMEHeaderKey(selectors[1]))

		switch selectors[0] {
		case "form":
			cfg.Extractor = CsrfFromForm(selectors[1])
		case "query":
			cfg.Extractor = CsrfFromQuery(selectors[1])
		case "param":
			cfg.Extractor = CsrfFromParam(selectors[1])
		case "cookie":
			cfg.Extractor = CsrfFromCookie(selectors[1])
		}
	}

	return cfg
}

var (
	errTokenNotFound = errors.New("csrf token not found")
)

// Csrf creates a new middleware handler
func Csrf(config ...ConfigCsrf) http.HandlerFunc {
	// Set default config
	cfg := configDefault(config...)

	// Create manager to simplify storage operations ( see manager.go )
	manager := newManager(cfg.Storage)

	dummyValue := []byte{'+'}

	// Return new handler
	return func(c http.Context) (err error) {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		var token string

		// Action depends on the HTTP method
		switch c.Method() {
		case fiber.MethodGet, fiber.MethodHead, fiber.MethodOptions, fiber.MethodTrace:
			// Declare empty token and try to get existing CSRF from cookie
			token = c.Cookies(cfg.CookieName)
		default:
			// Assume that anything not defined as 'safe' by RFC7231 needs protection

			// Extract token from client request i.e. header, query, param, form or cookie
			token, err = cfg.Extractor(c)
			if err != nil {
				return cfg.ErrorHandler(c, err)
			}

			// if token does not exist in Storage
			if manager.getRaw(token) == nil {
				// Expire cookie
				c.Cookie(&http.Cookie{
					Name:        cfg.CookieName,
					Domain:      cfg.CookieDomain,
					Path:        cfg.CookiePath,
					Expires:     time.Now().Add(-1 * time.Minute),
					Secure:      cfg.CookieSecure,
					HTTPOnly:    cfg.CookieHTTPOnly,
					SameSite:    cfg.CookieSameSite,
					SessionOnly: cfg.CookieSessionOnly,
				})
				return cfg.ErrorHandler(c, errTokenNotFound)
			}
		}
		// Generate CSRF token if not exist
		if token == "" {
			// And generate a new token
			token = cfg.KeyGenerator()
		}

		// Add/update token to Storage
		manager.setRaw(token, dummyValue, cfg.Expiration)

		// Create cookie to pass token to client
		cookie := &http.Cookie{
			Name:        cfg.CookieName,
			Value:       token,
			Domain:      cfg.CookieDomain,
			Path:        cfg.CookiePath,
			Expires:     time.Now().Add(cfg.Expiration),
			Secure:      cfg.CookieSecure,
			HTTPOnly:    cfg.CookieHTTPOnly,
			SameSite:    cfg.CookieSameSite,
			SessionOnly: cfg.CookieSessionOnly,
		}
		// Set cookie to response
		c.Cookie(cookie)

		// Protect clients from caching the response by telling the browser
		// a new header value is generated
		c.Vary(fiber.HeaderCookie)

		// Store token in context if set
		if cfg.ContextKey != "" {
			c.WithValue(cfg.ContextKey, token)
		}

		// Continue stack
		return c.Next()
	}
}

var (
	errMissingHeader = errors.New("missing csrf token in header")
	errMissingQuery  = errors.New("missing csrf token in query")
	errMissingParam  = errors.New("missing csrf token in param")
	errMissingForm   = errors.New("missing csrf token in form")
	errMissingCookie = errors.New("missing csrf token in cookie")
)

// CsrfFromParam returns a function that extracts token from the url param string.
func CsrfFromParam(param string) func(c http.Context) (string, error) {
	return func(c http.Context) (string, error) {
		token := c.Params(param)
		if token == "" {
			return "", errMissingParam
		}
		return token, nil
	}
}

// CsrfFromForm returns a function that extracts a token from a multipart-form.
func CsrfFromForm(param string) func(c http.Context) (string, error) {
	return func(c http.Context) (string, error) {
		token := c.Form(param, "")
		if token == "" {
			return "", errMissingForm
		}
		return token, nil
	}
}

// CsrfFromCookie returns a function that extracts token from the cookie header.
func CsrfFromCookie(param string) func(c http.Context) (string, error) {
	return func(c http.Context) (string, error) {
		token := c.Cookies(param)
		if token == "" {
			return "", errMissingCookie
		}
		return token, nil
	}
}

// CsrfFromHeader returns a function that extracts token from the request header.
func CsrfFromHeader(param string) func(c http.Context) (string, error) {
	return func(c http.Context) (string, error) {
		token := c.Header(param, "")
		if token == "" {
			return "", errMissingHeader
		}
		return token, nil
	}
}

// CsrfFromQuery returns a function that extracts token from the query string.
func CsrfFromQuery(param string) func(c http.Context) (string, error) {
	return func(c http.Context) (string, error) {
		token := c.Query(param, "")
		if token == "" {
			return "", errMissingQuery
		}
		return token, nil
	}
}

// go:generate msgp
// msgp -file="manager.go" -o="manager_msgp.go" -tests=false -unexported
// don't forget to replace the msgp import path to:
// "github.com/gofiber/fiber/v2/internal/msgp"
type item struct{}

//msgp:ignore manager
type manager struct {
	pool    sync.Pool
	memory  *storage.Storage
	storage fiber.Storage
}

func newManager(str fiber.Storage) *manager {
	// Create new storage handler
	manager := &manager{
		pool: sync.Pool{
			New: func() interface{} {
				return new(item)
			},
		},
	}
	if str != nil {
		// Use provided storage if provided
		manager.storage = str
	} else {
		// Fallback too memory storage
		manager.memory = storage.New()
	}
	return manager
}

// acquire returns an *entry from the sync.Pool
func (m *manager) acquire() *item {
	return m.pool.Get().(*item)
}

// release and reset *entry to sync.Pool
func (m *manager) release(e *item) {
	// don't release item if we using memory storage
	if m.storage != nil {
		return
	}
	m.pool.Put(e)
}

// get data from storage or memory
func (m *manager) get(key string) (it *item) {
	if m.storage != nil {
		it = m.acquire()
		if raw, _ := m.storage.Get(key); raw != nil {
			if _, err := it.UnmarshalMsg(raw); err != nil {
				return
			}
		}
		return
	}
	if it, _ = m.memory.Get(key).(*item); it == nil {
		it = m.acquire()
	}
	return
}

// get raw data from storage or memory
func (m *manager) getRaw(key string) (raw []byte) {
	if m.storage != nil {
		raw, _ = m.storage.Get(key)
	} else {
		raw, _ = m.memory.Get(key).([]byte)
	}
	return
}

// set data to storage or memory
func (m *manager) set(key string, it *item, exp time.Duration) {
	if m.storage != nil {
		if raw, err := it.MarshalMsg(nil); err == nil {
			_ = m.storage.Set(key, raw, exp)
		}
	} else {
		// the key is crucial in crsf and sometimes a reference to another value which can be reused later(pool/unsafe values concept), so a copy is made here
		m.memory.Set(utils.CopyString(key), it, exp)
	}
}

// set data to storage or memory
func (m *manager) setRaw(key string, raw []byte, exp time.Duration) {
	if m.storage != nil {
		_ = m.storage.Set(key, raw, exp)
	} else {
		// the key is crucial in crsf and sometimes a reference to another value which can be reused later(pool/unsafe values concept), so a copy is made here
		m.memory.Set(utils.CopyString(key), raw, exp)
	}
}

// delete data from storage or memory
func (m *manager) delete(key string) {
	if m.storage != nil {
		_ = m.storage.Delete(key)
	} else {
		m.memory.Delete(key)
	}
}

// DecodeMsg implements msgp.Decodable
func (z *item) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z item) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 0
	err = en.Append(0x80)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z item) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 0
	o = append(o, 0x80)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *item) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z item) Msgsize() (s int) {
	s = 1
	return
}
