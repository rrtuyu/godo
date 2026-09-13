package main

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"godo/internal/db"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context

	db *sqlx.DB
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.db = db.InitDb(a.ctx)
	wailsruntime.LogInfo(a.ctx, "startup: DB inittialized successfully")
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) shutdown(ctx context.Context) {
	a.db.Close()
}
