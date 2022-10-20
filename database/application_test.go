package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sujit-baniya/framework/config"
)

func TestInit(t *testing.T) {
	configApp := config.ServiceProvider{}
	configApp.Register()

	assert.NotPanics(t, func() {
		app := Application{}
		app.Init()
	})
}
