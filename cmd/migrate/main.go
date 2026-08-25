// cmd/migrate/main.go
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"project/internal/core"
	"project/migrations"
)

func main() {
	if err := run(); err != nil {
		log.Printf("migrate: %v", err)
		os.Exit(1)
	}
}

func run() error {
	command := flag.String("cmd", "up", "goose command: up, down, status, redo, reset, version")
	flag.Parse()

	cfg := core.LoadConfig()

	db, err := sql.Open("pgx", cfg.DataBaseConfig.GetURL())
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}

	if err := goose.RunContext(context.Background(), *command, db, "."); err != nil {
		return fmt.Errorf("goose %s: %w", *command, err)
	}
	return nil
}
