package database

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect() {
	var err error
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	DB, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	if err := DB.Ping(context.Background()); err != nil {
		log.Fatalf("unable to ping databese:%v\n", err)
	}

	log.Printf("successfully connected to databse\n")
}

func RunMigrations(migrationsDir string) {
	log.Println("Running migrations...")

	// Create tracking table if it doesn't exist
	_, err := DB.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create schema_migrations table: %v", err)
	}

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Printf("Warning: failed to read migrations directory: %v", err)
		return
	}

	var upMigrations []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".up.sql") {
			upMigrations = append(upMigrations, f.Name())
		}
	}
	sort.Strings(upMigrations)

	for _, m := range upMigrations {
		var exists bool
		err := DB.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", m).Scan(&exists)
		if err != nil {
			log.Fatalf("Failed to check migration status for %s: %v", m, err)
		}

		if exists {
			log.Printf("Skipping migration %s (already applied)", m)
			continue
		}

		path := filepath.Join(migrationsDir, m)
		content, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("Failed to read migration %s: %v", m, err)
		}

		log.Printf("Executing migration: %s", m)

		// Execute migration and record it in a transaction
		tx, err := DB.Begin(context.Background())
		if err != nil {
			log.Fatalf("Failed to start transaction for %s: %v", m, err)
		}

		_, err = tx.Exec(context.Background(), string(content))
		if err != nil {
			tx.Rollback(context.Background())
			log.Fatalf("Failed to execute migration %s: %v", m, err)
		}

		_, err = tx.Exec(context.Background(), "INSERT INTO schema_migrations (version) VALUES ($1)", m)
		if err != nil {
			tx.Rollback(context.Background())
			log.Fatalf("Failed to record migration %s: %v", m, err)
		}

		err = tx.Commit(context.Background())
		if err != nil {
			log.Fatalf("Failed to commit migration %s: %v", m, err)
		}
	}
	log.Println("Migrations applied successfully.")
}
