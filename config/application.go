package config

import (
	"os"
	"slices"

	"github.com/gookit/color"
	"github.com/spf13/cast"
	"github.com/spf13/viper"

	"github.com/oarkflow/framework/contracts/config"
	"github.com/oarkflow/framework/facades"
	"github.com/oarkflow/framework/support"
	"github.com/oarkflow/framework/support/file"
)

func init() {
	if slices.Contains(os.Args, "artisan") {
		support.Env = support.EnvArtisan
	}
}

type Application struct {
	vip *viper.Viper
}

func Init() {
	app := Application{}
	facades.Config = app.Init()
}

func (app *Application) Init() config.Config {
	envFile := ".env"
	if !file.Exists(envFile) {
		color.Redln("Please create .env and initialize it first.")
		color.Warnln("Run command: \ncp .env.example .env && go run . artisan key:generate")
		os.Exit(0)
	}

	app.vip = viper.New()
	app.vip.SetConfigName(envFile)
	app.vip.SetConfigType("env")
	app.vip.AddConfigPath(".")
	err := app.vip.ReadInConfig()
	if err != nil {
		panic(err.Error())
	}
	app.vip.SetEnvPrefix("goravel")
	app.vip.AutomaticEnv()
	appKey := app.Env("APP_KEY")
	if appKey == nil && support.Env != support.EnvArtisan {
		color.Redln("Please initialize APP_KEY first.")
		color.Warnln("Run command: \ngo run . artisan key:generate")
		os.Exit(0)
	} else if appKey == nil {
		return app
	}

	if len(appKey.(string)) > 0 && len(appKey.(string)) != 32 {
		color.Redln("Invalid APP_KEY, please reset it.")
		color.Warnln("Run command: \ngo run . artisan key:generate")
		os.Exit(0)
	}
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
