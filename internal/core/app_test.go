package core

import (
	"path/filepath"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
)

func TestNewApp(t *testing.T) {
	path := filepath.Join("testdata", "config", "valid.yaml")
	app, err := NewApp("", "", "", path)
	defer app.Close()

	assert.NotNil(t, app)
	assert.NoError(t, err)

	assert.Equal(t, path, app.Config.ConfigPath)
}

func TestNewApp_WhenConfigError(t *testing.T) {
	path := filepath.Join("testdata", "config", "malformed.yaml")
	app, err := NewApp("", "", "", path)

	assert.Nil(t, app)
	assert.Error(t, err)
}

func TestNewTestApp(t *testing.T) {
	app := NewTestApp()
	defer app.Close()

	assert.NotNil(t, app)
}

func TestAppForContext(t *testing.T) {
	app := NewTestApp()
	ctx := app.Context()

	assert.Equal(t, app, AppForContext(ctx))
}

func TestApp_SetVerbosity(t *testing.T) {
	app := NewTestApp()

	// default
	assert.Equal(t, log.WarnLevel, app.Logger.GetLevel())

	app.SetVerbosity(0) // noop
	assert.Equal(t, log.WarnLevel, app.Logger.GetLevel())
	app.SetVerbosity(1) // info
	assert.Equal(t, log.InfoLevel, app.Logger.GetLevel())
	app.SetVerbosity(2) // debug
	assert.Equal(t, log.DebugLevel, app.Logger.GetLevel())
	app.SetVerbosity(99) // capped to debug
	assert.Equal(t, log.DebugLevel, app.Logger.GetLevel())
}
