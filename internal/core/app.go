package core

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/twelvelabs/termite/ui"
)

type ctxKey string

var (
	ctxKeyApp ctxKey = "github.com/twelvelabs/envctl/internal/core.App"
)

// App contains global and/or singleton application data.
type App struct {
	Config *Config
	Meta   *Meta
	IO     *ui.IOStreams
	UI     *ui.UserInterface

	ctx context.Context //nolint: containedctx
}

// AppForContext returns the app singleton stored in the given context.
func AppForContext(ctx context.Context) *App {
	return ctx.Value(ctxKeyApp).(*App)
}

// NewApp returns the default App singleton.
func NewApp(version, commit, date, path string) (*App, error) {
	config, err := NewConfigFromPath(path)
	if err != nil {
		return nil, err
	}

	meta := NewMeta(version, commit, date)
	ios := ui.NewIOStreams()

	app := &App{
		Config: config,
		Meta:   meta,
		IO:     ios,
		UI:     ui.NewUserInterface(ios),
	}

	return app, nil
}

// NewTestApp returns the test App singleton.
// All properties will be configured for testing (mocks, stubs, etc).
func NewTestApp() *App {
	config, _ := NewTestConfig()

	meta := NewMeta("test", "", "0")
	ios := ui.NewTestIOStreams()

	app := &App{
		Config: config,
		Meta:   meta,
		IO:     ios,
		UI:     ui.NewUserInterface(ios),
	}

	return app
}

// Close ensures all app resources have been closed.
func (a *App) Close() error {
	return nil
}

// Context returns the root [context.Context] for the app.
// The root context is automatically configured for graceful termination,
// and will be canceled on SIGINT or SIGTERM.
//
// Application logic should make use of the context done channel:
//
//	// do long running logic in goroutine and block on done channel
//	go func() { ... }()
//	<-ctx.Done()
//	// when done channel closes (i.e. SIGINT received), cancel operation
//	cancelOperation()
func (a *App) Context() context.Context {
	if a.ctx == nil {
		done := make(chan os.Signal, 1)
		signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			sig := <-done
			fmt.Fprintf(a.IO.Err, "Received signal: %+v\n", sig)
			cancel()
		}()

		a.ctx = context.WithValue(ctx, ctxKeyApp, a)
	}
	return a.ctx
}
