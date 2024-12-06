package main

import (
	"context"
	"log/slog"
	"os"
	"sync"

	"github.com/regalias/scry/pkg/buylist"
	"github.com/regalias/scry/pkg/vendors"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx      context.Context
	logger   *slog.Logger
	vendors  *vendors.Manager
	buylists *buylist.Manager
	lock     sync.RWMutex
	ready    bool
	hasError bool
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {

	a.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	a.ctx = ctx
	a.lock.Lock()
	defer a.lock.Unlock()
	a.ready = false
}

func (a *App) onDomReady(ctx context.Context) {

	runtime.WindowSetDarkTheme(ctx)

	vendors, err := vendors.NewManager(a.logger)
	if err != nil {
		a.logger.Error("failed to create vendors manager", "err", err)
		runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Error",
			Message: err.Error(),
		})
		a.lock.Lock()
		defer a.lock.Unlock()
		a.hasError = true
		return
	}

	buylists, err := buylist.NewManager(a.logger)
	if err != nil {
		a.logger.Error("failed to create buylist manager", "err", err)
		runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Error",
			Message: err.Error(),
		})
		a.lock.Lock()
		defer a.lock.Unlock()
		a.hasError = true
		return
	}

	a.vendors = vendors
	a.buylists = buylists
	a.lock.Lock()
	defer a.lock.Unlock()
	a.ready = true
}

func (a *App) IsReady() bool {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.ready
}

func (a *App) IsError() bool {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.hasError
}

func (a *App) shutdown(_ context.Context) {
	a.buylists.Shutdown()
}
