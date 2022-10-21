package console

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sujit-baniya/framework/config"
	"github.com/sujit-baniya/framework/console"
	"github.com/sujit-baniya/framework/contracts"
	console2 "github.com/sujit-baniya/framework/contracts/console"
	"github.com/sujit-baniya/framework/facades"
	"github.com/sujit-baniya/framework/support/file"
	testingfile "github.com/sujit-baniya/framework/testing/file"
)

func TestListenerMakeCommand(t *testing.T) {
	err := testingfile.CreateEnv()
	assert.Nil(t, err)

	configApp := config.ServiceProvider{}
	configApp.Register()

	facadesConfig := facades.Config
	facadesConfig.Add("app", map[string]any{
		"providers": []contracts.ServiceProvider{},
	})

	consoleApp := console.Application{}
	instance := consoleApp.Init()
	instance.Register([]console2.Command{
		&ListenerMakeCommand{},
	})

	assert.NotPanics(t, func() {
		instance.Call("make:listener GoravelListen")
	})

	assert.True(t, file.Exist("app/listeners/goravel_listen.go"))
	assert.True(t, file.Remove("app"))
	err = os.Remove(".env")
	assert.Nil(t, err)
}
