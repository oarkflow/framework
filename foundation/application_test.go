package foundation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sujit-baniya/framework/config"
	"github.com/sujit-baniya/framework/console"
	"github.com/sujit-baniya/framework/contracts"
	"github.com/sujit-baniya/framework/facades"
)

func TestInit(t *testing.T) {
	configApp := config.ServiceProvider{}
	configApp.Register()

	facadesConfig := facades.Config
	facadesConfig.Add("app", map[string]any{
		"providers": []contracts.ServiceProvider{
			&console.ServiceProvider{},
		},
	})

	assert.NotPanics(t, func() {
		app := Application{}
		app.Boot()
	})
}
