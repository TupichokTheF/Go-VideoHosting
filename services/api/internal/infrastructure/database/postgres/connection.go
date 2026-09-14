package postgres

import (
	"context"
	"fmt"
	"log"
	"project/internal/core"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Connection struct {
	ConnPool *pgxpool.Pool
}

func OpenConnection(config core.DataBaseConfig) (*Connection, error) {
	db, err := pgxpool.New(context.Background(), config.GetURL())
	if err != nil {
		return nil, fmt.Errorf("Error while oppening database connection: %s", err)
	}

	if err := db.Ping(context.TODO()); err != nil {
		log.Fatalf("ping: %v", err)
	}

	return &Connection{ConnPool: db}, nil
}

func (c *Connection) CloseConnection() {
	c.ConnPool.Close()
}
