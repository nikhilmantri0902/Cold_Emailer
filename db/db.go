package db

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

const (
	MigrationsDir  = "migrations"
	MigrationTable = "migration_metadata"
)

// InitDB initializes the global DB instance
func InitDB(connStr string) error {
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Println("error:", err)
		return err
	}
	if err := DB.Ping(); err != nil {
		log.Println("error:", err)
		return err
	}
	return nil
}

// GetDB returns the global DB instance
func GetDB() *sql.DB {
	return DB
}

// EnsureMigrationTable ensures the migration_metadata table exists
func EnsureMigrationTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS ` + MigrationTable + ` (
		id SERIAL PRIMARY KEY,
		filename TEXT NOT NULL UNIQUE,
		run_at TIMESTAMP NOT NULL DEFAULT NOW()
	)`)
	if err != nil {
		log.Println("error:", err)
	}
	return err
}

// GetAppliedMigrations returns a sorted list of applied migration filenames
func GetAppliedMigrations(db *sql.DB) ([]string, error) {
	rows, err := db.Query(`SELECT filename FROM ` + MigrationTable + ` ORDER BY filename ASC`)
	if err != nil {
		log.Println("error:", err)
		return nil, err
	}
	defer rows.Close()
	var applied []string
	for rows.Next() {
		var fname string
		if err := rows.Scan(&fname); err != nil {
			log.Println("error:", err)
			return nil, err
		}
		applied = append(applied, fname)
	}
	return applied, nil
}

// ListMigrationFiles returns a sorted list of migration files in the migrations directory
func ListMigrationFiles() ([]string, error) {
	files, err := ioutil.ReadDir(MigrationsDir)
	if err != nil {
		log.Println("error:", err)
		return nil, err
	}
	var migrationFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
			migrationFiles = append(migrationFiles, f.Name())
		}
	}
	sort.Strings(migrationFiles)
	return migrationFiles, nil
}

// RunMigrations runs all pending migrations in order
func RunMigrations(db *sql.DB) error {
	if err := EnsureMigrationTable(db); err != nil {
		log.Println("error:", err)
		return err
	}
	applied, err := GetAppliedMigrations(db)
	if err != nil {
		log.Println("error:", err)
		return err
	}
	appliedSet := make(map[string]struct{})
	for _, fname := range applied {
		appliedSet[fname] = struct{}{}
	}
	files, err := ListMigrationFiles()
	if err != nil {
		log.Println("error:", err)
		return err
	}
	for _, fname := range files {
		if _, already := appliedSet[fname]; already {
			continue
		}
		log.Printf("Running migration: %s", fname)
		content, err := ioutil.ReadFile(filepath.Join(MigrationsDir, fname))
		if err != nil {
			log.Println("error:", err)
			return err
		}
		if _, err := db.Exec(string(content)); err != nil {
			log.Println("error:", err)
			return fmt.Errorf("migration %s failed: %w", fname, err)
		}
		if _, err := db.Exec(`INSERT INTO `+MigrationTable+` (filename) VALUES ($1)`, fname); err != nil {
			log.Println("error:", err)
			return err
		}
	}
	return nil
}

// GenerateMigrationFile creates a new migration file with timestamp and suffix
func GenerateMigrationFile(suffix string) (string, error) {
	if suffix == "" {
		log.Println("error: suffix required")
		return "", errors.New("suffix required")
	}
	if err := os.MkdirAll(MigrationsDir, 0755); err != nil {
		log.Println("error:", err)
		return "", err
	}
	timestamp := time.Now().Format("20060102_150405")
	fname := fmt.Sprintf("%s_%s.sql", timestamp, suffix)
	path := filepath.Join(MigrationsDir, fname)
	if err := ioutil.WriteFile(path, []byte("-- SQL migration\n"), 0644); err != nil {
		log.Println("error:", err)
		return "", err
	}
	return path, nil
}
