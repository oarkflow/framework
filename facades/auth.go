package facades

import (
	"github.com/oarkflow/framework/auth/access"
	"github.com/oarkflow/framework/contracts/auth"
)

var (
	Gate        *access.Gate
	Auth        auth.Auth
	JwtAuth     auth.Auth
	ApiKeyAuth  auth.Auth
	SessionAuth auth.Auth
)
