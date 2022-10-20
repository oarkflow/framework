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

func TestJobMakeCommand(t *testing.T) {
	err := testingfile.CreateEnv()
	assert.Nil(t, err)

	configApp := config.ServiceProvider{}
	configApp.Register()

	facadesConfig := facades.Config
	facadesConfig.Add("app", map[string]interface{}{
		"providers": []contracts.ServiceProvider{},
	})

	consoleApp := console.Application{}
	instance := consoleApp.Init()
	instance.Register([]console2.Command{
		&JobMakeCommand{},
	})

	assert.NotPanics(t, func() {
		instance.Call("make:job GoravelJob")
	})

	assert.True(t, file.Exist("app/jobs/goravel_job.go"))
	assert.True(t, file.Remove("app"))
	err = os.Remove(".env")
	assert.Nil(t, err)
}
