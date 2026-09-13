package db

import (
	"context"
	"database/sql"
	"path/filepath"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
	_ "modernc.org/sqlite"
)

func InitDb(ctx context.Context) *sql.DB {
	db, err := sql.Open("sqlite", "test.db")
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	
	migrator := newMigrator(db)
	if err := migrator.Run(); err != nil {
		panic(err)
	}

	if absPath, err := filepath.Abs("test.db"); err == nil {
		wailsruntime.LogInfo(ctx, "DB path: "+absPath)
	}
	return db
}