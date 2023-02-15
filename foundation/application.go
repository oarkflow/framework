package foundation

import (
	"github.com/sujit-baniya/framework/console"
	"os"
	"strings"

	"github.com/sujit-baniya/framework/config"
	"github.com/sujit-baniya/framework/contracts"
	"github.com/sujit-baniya/framework/facades"
	"github.com/sujit-baniya/framework/support"
)

func Init(providers ...contracts.ServiceProvider) *Application {
	//Create a new application instance.
	app := &Application{Providers: providers}
	app.registerBaseServiceProviders()
	app.bootBaseServiceProviders()
	return app
}

type Application struct {
	Providers []contracts.ServiceProvider
}

// Boot Register and bootstrap configured service providers.
func (app *Application) Boot() {
	app.registerConfiguredServiceProviders()
	app.bootConfiguredServiceProviders()

	app.bootArtisan()
	app.setRootPath()
}

func (app *Application) setRootPath() {
	rootPath := getCurrentAbPath()

	// Hack air path
	airPath := "/storage/temp"
	if strings.HasSuffix(rootPath, airPath) {
		rootPath = strings.ReplaceAll(rootPath, airPath, "")
	}

	support.RootPath = rootPath
}

// bootArtisan Boot artisan command.
func (app *Application) bootArtisan() {
	facades.Artisan.Run(os.Args, true)
}

// getBaseServiceProviders Get base service providers.
func (app *Application) getBaseServiceProviders() []contracts.ServiceProvider {
	return []contracts.ServiceProvider{
		&config.ServiceProvider{},
		&console.ServiceProvider{},
	}
}

// getConfiguredServiceProviders Get configured service providers.
func (app *Application) getConfiguredServiceProviders() []contracts.ServiceProvider {
	return app.Providers
}

// registerBaseServiceProviders Register base service providers.
func (app *Application) registerBaseServiceProviders() {
	app.registerServiceProviders(app.getBaseServiceProviders())
}

// bootBaseServiceProviders Bootstrap base service providers.
func (app *Application) bootBaseServiceProviders() {
	app.bootServiceProviders(app.getBaseServiceProviders())
}

// registerConfiguredServiceProviders Register configured service providers.
func (app *Application) registerConfiguredServiceProviders() {
	app.registerServiceProviders(app.getConfiguredServiceProviders())
}

// bootConfiguredServiceProviders Bootstrap configured service providers.
func (app *Application) bootConfiguredServiceProviders() {
	app.bootServiceProviders(app.getConfiguredServiceProviders())
}

// registerServiceProviders Register service providers.
func (app *Application) registerServiceProviders(serviceProviders []contracts.ServiceProvider) {
	for _, serviceProvider := range serviceProviders {
		serviceProvider.Register()
	}
}

// bootServiceProviders Bootstrap service providers.
func (app *Application) bootServiceProviders(serviceProviders []contracts.ServiceProvider) {
	for _, serviceProvider := range serviceProviders {
		serviceProvider.Boot()
	}
}

// RunningInConsole Determine if the application is running in the console.
func (app *Application) RunningInConsole() bool {
	args := os.Args

	return len(args) >= 2 && args[1] == "artisan"
}
