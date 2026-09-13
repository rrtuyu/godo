package db

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"sort"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

const migrationDir = "migrations"
const migrationTable = "schema_migrations"

type Migrator struct {
	db *sql.DB
}

func newMigrator(db *sql.DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) Run() error {
	query := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s 
	(
		version VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`, migrationTable)
	_, err := m.db.Exec(query)
	if err != nil {
		return err
	}

	files, err := fs.ReadDir(migrationFiles, migrationDir)
	if err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		version := file.Name()
		var exists bool
		err := m.db.QueryRow(fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE version = ?)", migrationTable), version).Scan(&exists)
		if err != nil {
			return err
		}

		if exists {
			continue
		}

		content, err := migrationFiles.ReadFile(migrationDir + "/" + version)
		if err != nil {
			return err
		}

		err = m.inTransaction(func(tx *sql.Tx) error {
			if _, err := tx.Exec(string(content)); err != nil {
				return fmt.Errorf("execute SQL error: %w", err)
			}

			if _, err := tx.Exec("INSERT INTO "+migrationTable+" (version) values ($1)", version); err != nil {
				return fmt.Errorf("failed to log migration: %w", err)
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Migrator) inTransaction(fn func(tx *sql.Tx) error) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
