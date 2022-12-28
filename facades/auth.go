package facades

import (
	"github.com/sujit-baniya/framework/auth/access"
	"github.com/sujit-baniya/framework/contracts/auth"
)

var (
	Gate        *access.Gate
	Auth        auth.Auth
	JwtAuth     auth.Auth
	ApiKeyAuth  auth.Auth
	SessionAuth auth.Auth
)
