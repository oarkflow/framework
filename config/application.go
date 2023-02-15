package config

import (
	"github.com/gookit/color"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"github.com/sujit-baniya/framework/support/file"
	"os"
)

type Application struct {
	vip *viper.Viper
}

func NewApplication(envPath string) *Application {
	if !file.Exists(envPath) {
		color.Redln("Please create .env and initialize it first\nRun command: \ncp .env.example .env && go run . artisan key:generate")
		os.Exit(0)
	}

	app := &Application{}
	app.vip = viper.New()
	app.vip.SetConfigName(envPath)
	app.vip.SetConfigType("env")
	app.vip.AddConfigPath(".")
	err := app.vip.ReadInConfig()
	if err != nil {
		panic(err.Error())
	}
	app.vip.SetEnvPrefix("goravel")
	app.vip.AutomaticEnv()

	return app
}

// Env Get config from env.
func (app *Application) Env(envName string, defaultValue ...any) any {
	value := app.Get(envName, defaultValue...)
	if cast.ToString(value) == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}

		return nil
	}

	return value
}

// Add config to application.
func (app *Application) Add(name string, configuration map[string]any) {
	app.vip.Set(name, configuration)
}

// Get config from application.
func (app *Application) Get(path string, defaultValue ...any) any {
	if !app.vip.IsSet(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}

	return app.vip.Get(path)
}

// GetString Get string type config from application.
func (app *Application) GetString(path string, defaultValue ...any) string {
	value := cast.ToString(app.Get(path, defaultValue...))
	if value == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0].(string)
		}

		return ""
	}

	return value
}

// GetInt Get int type config from application.
func (app *Application) GetInt(path string, defaultValue ...any) int {
	value := app.Get(path, defaultValue...)
	if cast.ToString(value) == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0].(int)
		}

		return 0
	}

	return cast.ToInt(value)
}

// GetBool Get bool type config from application.
func (app *Application) GetBool(path string, defaultValue ...any) bool {
	value := app.Get(path, defaultValue...)
	if cast.ToString(value) == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0].(bool)
		}

		return false
	}

	return cast.ToBool(value)
}
