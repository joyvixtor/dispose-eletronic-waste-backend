package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DB *sql.DB
}

func New(dbPath string) (*Database, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &Database{DB: db}

	if err := database.migrate(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database initialized successfully")
	return database, nil
}

// migrate applies all pending .up.sql migrations found in the migrations directory.
// It creates a schema_migrations table to keep track of applied migrations.
func (d *Database) migrate() error {
	// Allow overriding migrations dir via env, default to ./migrations
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = "migrations"
	}

	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		// No migrations directory; nothing to do
		log.Printf("migrations directory '%s' not found, skipping migrations", migrationsDir)
		return nil
	}

	if _, err := d.DB.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY, applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP)`); err != nil {
		return fmt.Errorf("failed to ensure schema_migrations table: %w", err)
	}

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations dir: %w", err)
	}

	// Collect .up.sql files
	type migration struct {
		version string
		path    string
	}
	var migrations []migration
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() {
			continue
		}
		if filepath.Ext(name) != ".sql" {
			continue
		}
		if !strings.HasSuffix(name, ".up.sql") {
			continue
		}
		version := strings.TrimSuffix(strings.TrimSuffix(name, ".sql"), ".up") // remove .up.sql
		migrations = append(migrations, migration{version: version, path: filepath.Join(migrationsDir, name)})
	}

	sort.Slice(migrations, func(i, j int) bool { return migrations[i].version < migrations[j].version })

	tx, err := d.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin migration tx: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	for _, m := range migrations {
		var exists string
		err = tx.QueryRow(`SELECT version FROM schema_migrations WHERE version = ?`, m.version).Scan(&exists)
		if err != nil && err != sql.ErrNoRows {
			_ = tx.Rollback()
			return fmt.Errorf("failed to query schema_migrations: %w", err)
		}
		if exists != "" {
			continue
		} // already applied

		sqlBytes, err := os.ReadFile(m.path)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to read migration %s: %w", m.path, err)
		}

		log.Printf("applying migration %s", m.version)
		if _, err := tx.Exec(string(sqlBytes)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", m.version, err)
		}

		if _, err := tx.Exec(`INSERT INTO schema_migrations (version) VALUES (?)`, m.version); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", m.version, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migrations: %w", err)
	}
	return nil
}
