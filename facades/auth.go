package facades

import (
	"github.com/sujit-baniya/framework/contracts/auth"
)

var (
	Auth        auth.Auth
	JwtAuth     auth.Auth
	ApiKeyAuth  auth.Auth
	SessionAuth auth.Auth
)
