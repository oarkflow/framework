package cache

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/sujit-baniya/framework/config"
	"github.com/sujit-baniya/framework/console"
	"github.com/sujit-baniya/framework/facades"
	"github.com/sujit-baniya/framework/testing/file"
)

func TestInit(t *testing.T) {
	initConfig()

	assert.NotPanics(t, func() {
		app := Application{}
		app.Init()
	})
}

func TestClearCommand(t *testing.T) {
	initConfig()

	consoleServiceProvider := console.ServiceProvider{}
	consoleServiceProvider.Register()

	cacheServiceProvider := ServiceProvider{}
	cacheServiceProvider.Register()
	cacheServiceProvider.Boot()

	err := facades.Cache.Put("test-clear-command", "goravel", 5*time.Second)
	assert.Nil(t, err)
	assert.True(t, facades.Cache.Has("test-clear-command"))

	assert.NotPanics(t, func() {
		facades.Artisan.Call("cache:clear")
	})

	assert.False(t, facades.Cache.Has("test-clear-command"))
}

func initConfig() {
	file.CreateEnv()
	configServiceProvider := config.ServiceProvider{}
	configServiceProvider.Register()

	facadesConfig := facades.Config
	facadesConfig.Add("cache", map[string]any{
		"default": facadesConfig.Env("CACHE_DRIVER", "redis"),
		"stores": map[string]any{
			"redis": map[string]any{
				"driver":     "redis",
				"connection": "default",
			},
		},
		"prefix": "goravel_cache",
	})

	facadesConfig.Add("database", map[string]any{
		"redis": map[string]any{
			"default": map[string]any{
				"host":     facadesConfig.Env("REDIS_HOST", "127.0.0.1"),
				"password": facadesConfig.Env("REDIS_PASSWORD", ""),
				"port":     facadesConfig.Env("REDIS_PORT", 6379),
				"database": facadesConfig.Env("REDIS_DB", 0),
			},
		},
	})

	os.Remove(".env")
}
