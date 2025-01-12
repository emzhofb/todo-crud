package database

import (
	"context"
	"fmt"

	"github.com/emzhofb/todo-crud/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgreSQL Implementation
type PostgresDB struct {
	connString string
}

func NewPostgresDB(config util.Config) *PostgresDB {
	return &PostgresDB{
		connString: config.DBSource,
	}
}

func (p *PostgresDB) Connect() error {
	_, err := pgxpool.New(context.Background(), p.connString)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDB) Close() error {
	fmt.Println("Closing PostgreSQL connection")
	// Add logic to close PostgreSQL connection
	return nil
}

func (p *PostgresDB) Ping() error {
	fmt.Println("Pinging PostgreSQL")
	// Add logic to ping PostgreSQL
	return nil
}
