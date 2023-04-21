package envctl

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/twelvelabs/termite/ui"
)

type ctxKey string

var (
	ctxKeyApp ctxKey = "github.com/twelvelabs/envctl/internal/envctl.App"
)

// AppForContext returns the app singleton stored in the given context.
func AppForContext(ctx context.Context) *App {
	return ctx.Value(ctxKeyApp).(*App)
}

// NewApp returns the default App singleton.
func NewApp(version, commit, date string) (*App, error) {
	meta := NewAppMeta(version, commit, date)
	ios := ui.NewIOStreams()

	app := &App{
		Meta: meta,
		IO:   ios,
		UI:   ui.NewUserInterface(ios),
	}

	return app, nil
}

// NewTestApp returns the test App singleton.
// All properties will be configured for testing (mocks, stubs, etc).
func NewTestApp() *App {
	meta := NewAppMeta("test", "", "0")
	ios := ui.NewTestIOStreams()

	app := &App{
		Meta: meta,
		IO:   ios,
		UI:   ui.NewUserInterface(ios),
	}

	return app
}

// App contains global and/or singleton application data.
type App struct {
	Meta *AppMeta
	IO   *ui.IOStreams
	UI   *ui.UserInterface

	ctx context.Context //nolint: containedctx
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

// NewAppMeta returns a new AppMeta struct.
func NewAppMeta(version, commit, date string) *AppMeta {
	buildTime, _ := time.Parse(time.RFC3339, date)

	meta := &AppMeta{
		BuildCommit: commit,
		BuildTime:   buildTime,
		Version:     version,
		GOOS:        runtime.GOOS,
		GOARCH:      runtime.GOARCH,
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		meta.BuildGoVersion = info.GoVersion
		meta.BuildVersion = info.Main.Version
		meta.BuildChecksum = info.Main.Sum
	}

	return meta
}

// AppMeta contains application metadata (version, os, build info, etc).
type AppMeta struct {
	BuildCommit    string
	BuildTime      time.Time
	BuildGoVersion string
	BuildVersion   string
	BuildChecksum  string
	Version        string
	GOOS           string
	GOARCH         string
}
