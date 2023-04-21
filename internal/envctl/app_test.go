package envctl

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewApp(t *testing.T) {
	app, err := NewApp("", "", "")
	defer app.Close()

	assert.IsType(t, &App{}, app)
	assert.NotNil(t, app)
	assert.NoError(t, err)
}

func TestNewTestApp(t *testing.T) {
	app := NewTestApp()
	defer app.Close()

	assert.IsType(t, &App{}, app)
	assert.NotNil(t, app)
}

func TestAppForContext(t *testing.T) {
	app := NewTestApp()
	ctx := app.Context()
	assert.Equal(t, app, AppForContext(ctx))
}

func TestNewAppMeta(t *testing.T) {
	meta := NewAppMeta("1.2.3", "9b11774", "2023-02-19T00:57:51-06:00")
	assert.Equal(t, "9b11774", meta.BuildCommit)
	assert.Equal(t, "2023-02-19T00:57:51-06:00", meta.BuildTime.Format(time.RFC3339))
	assert.Equal(t, "1.2.3", meta.Version)
}
