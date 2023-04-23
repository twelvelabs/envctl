package core

import (
	"context"
	"time"

	"github.com/charmbracelet/log"
	"github.com/twelvelabs/termite/ui"
)

type ctxKey string

const (
	ctxKeyApp ctxKey = "github.com/twelvelabs/envctl/internal/core.App"

	verbosityInfo  = 1
	verbosityDebug = 2
)

// App contains global and/or singleton application data.
type App struct {
	CreatedAt time.Time
	Config    *Config
	Logger    *log.Logger
	Meta      *Meta
	IO        *ui.IOStreams
	UI        *ui.UserInterface

	ctx context.Context //nolint: containedctx
}

// AppForContext returns the app singleton stored in the given context.
func AppForContext(ctx context.Context) *App {
	return ctx.Value(ctxKeyApp).(*App)
}

// NewApp returns the default App singleton.
func NewApp(version, commit, date, path string) (*App, error) {
	created := time.Now()

	config, err := NewConfigFromPath(path)
	if err != nil {
		return nil, err
	}

	meta := NewMeta(version, commit, date)
	ios := ui.NewIOStreams()
	logger := newLogger(ios, config)

	app := &App{
		CreatedAt: created,
		Config:    config,
		Logger:    logger,
		Meta:      meta,
		IO:        ios,
		UI:        ui.NewUserInterface(ios),
	}

	return app, nil
}

// NewTestApp returns the test App singleton.
// All properties will be configured for testing (mocks, stubs, etc).
func NewTestApp() *App {
	created := time.Now()
	config, _ := NewTestConfig()

	meta := NewMeta("test", "", "0")
	ios := ui.NewTestIOStreams()
	logger := newLogger(ios, config)

	app := &App{
		CreatedAt: created,
		Config:    config,
		Logger:    logger,
		Meta:      meta,
		IO:        ios,
		UI:        ui.NewUserInterface(ios),
	}

	return app
}

// Close ensures all app resources have been closed.
func (a *App) Close() error {
	// Add any db/file closing here if needed.
	return nil
}

// Context returns the root [context.Context] for the app.
func (a *App) Context() context.Context {
	if a.ctx == nil {
		a.ctx = context.WithValue(context.Background(), ctxKeyApp, a)
	}
	return a.ctx
}

// SetVerbosity sets the log level for the given value.
//   - 1: INFO
//   - 2: DEBUG
func (a *App) SetVerbosity(value int) {
	if value > verbosityDebug {
		value = verbosityDebug
	}
	switch value {
	case verbosityInfo:
		a.Logger.SetLevel(log.InfoLevel)
	case verbosityDebug:
		a.Logger.SetLevel(log.DebugLevel)
	}
}

func newLogger(ios *ui.IOStreams, config *Config) *log.Logger {
	level := config.LogLevel
	if config.Debug {
		level = "debug"
	}
	return log.NewWithOptions(ios.Err, log.Options{
		Level:           log.ParseLevel(level),
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
	})
}
